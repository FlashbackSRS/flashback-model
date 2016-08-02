package fb

import (
	"encoding/json"
	"errors"
	"time"
)

type Note struct {
	ID
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	Model       *Model
	FieldValues []*FieldValue
	Attachments *FileCollection
}

/*
type Note struct {
    Tags           string           `db:"tags"` // List of the note's tags
    UniqueField    string           `db:"sfld"` // The text of the first field, used for Anki's simplistic uniqueness checking
    Checksum       int64            `db:"csum"` // Field checksum used for duplicate check. Integer representation of first 8 digits of sha1 hash of the first field
}
*/

type noteDoc struct {
	Type        string          `json:"type"`
	ID          ID              `json:"_id"`
	Rev         *string         `json:"_rev,omitempty"`
	Created     *time.Time      `json:"created,omitempty"`
	Modified    *time.Time      `json:"modified,omitempty"`
	Imported    *time.Time      `json:"imported,omitempty"`
	Model       string          `json:"model"`
	FieldValues []*FieldValue   `json:"fieldValues"`
	Attachments *FileCollection `json:"_attachments,omitempty"`
}

func NewNote(id string, model *Model) (*Note, error) {
	n := &Note{}
	nid, err := NewID("note", id)
	if err != nil {
		return nil, err
	}
	n.ID = nid
	n.Model = model
	n.FieldValues = make([]*FieldValue, len(model.Fields))
	n.Attachments = NewFileCollection()
	return n, nil
}

func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(noteDoc{
		Type:        "note",
		ID:          n.ID,
		Rev:         n.Rev,
		Created:     n.Created,
		Modified:    n.Modified,
		Imported:    n.Imported,
		Model:       n.Model.Identity(),
		FieldValues: n.FieldValues,
		Attachments: n.Attachments,
	})
}

func (n *Note) GetFieldValue(ord int) *FieldValue {
	fv := n.FieldValues[ord]
	if fv == nil {
		fv = &FieldValue{
			field: n.Model.Fields[ord],
		}
		n.FieldValues[ord] = fv
	}
	if fv.field.Type != TextField {
		fv.files = n.Attachments.NewView()
	}
	return fv
}

func (fv *FieldValue) Type() FieldType {
	return fv.field.Type
}

type FieldValue struct {
	field *Field
	text  *string
	files *FileCollectionView
}

type fieldValueDoc struct {
	Text  *string             `json:"text,omitempty"`
	Files *FileCollectionView `json:"files,omitempty"`
}

func (fv *FieldValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(fieldValueDoc{
		Text:  fv.text,
		Files: fv.files,
	})
}

func (fv *FieldValue) SetText(text string) error {
	if fv.field.Type != TextField && fv.field.Type != AnkiField {
		return errors.New("Text field not permitted")
	}
	fv.text = &text
	return nil
}

func (fv *FieldValue) AddFile(name, ctype string, content []byte) error {
	if fv.field.Type == TextField {
		return errors.New("Text fields do not support attachments")
	}
	return fv.files.AddFile(name, ctype, content)
}
