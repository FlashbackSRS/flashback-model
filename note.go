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
	ModelID     string
	FieldValues []*FieldValue
	Attachments *FileCollection
	model       *Model
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
	ModelID     string          `json:"model"`
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
	n.ModelID = model.Identity()
	n.FieldValues = make([]*FieldValue, len(model.Fields))
	n.Attachments = NewFileCollection()
	n.model = model
	return n, nil
}

func (n *Note) SetModel(m *Model) {
	n.model = m
	for i := 0; i < len(n.FieldValues); i++ {
		n.FieldValues[i].field = m.Fields[i]
	}
}

func (n *Note) Model() *Model {
	return n.model
}

func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(noteDoc{
		Type:        "note",
		ID:          n.ID,
		Rev:         n.Rev,
		Created:     n.Created,
		Modified:    n.Modified,
		Imported:    n.Imported,
		ModelID:     n.ModelID,
		FieldValues: n.FieldValues,
		Attachments: n.Attachments,
	})
}

func (n *Note) UnmarshalJSON(data []byte) error {
	doc := &noteDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "note" {
		return errors.New("Invalid document type for note: " + doc.Type)
	}
	n.ID = doc.ID
	n.Rev = doc.Rev
	n.Created = doc.Created
	n.Modified = doc.Modified
	n.Imported = doc.Imported
	n.ModelID = doc.ModelID
	n.FieldValues = doc.FieldValues
	n.Attachments = doc.Attachments
	for _, fv := range n.FieldValues {
		if fv.files != nil {
			n.Attachments.AddView(fv.files)
		}
	}
	return nil
}

func (n *Note) GetFieldValue(ord int) *FieldValue {
	fv := n.FieldValues[ord]
	if fv == nil {
		fv = &FieldValue{
			field: n.model.Fields[ord],
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

func (fv *FieldValue) UnmarshalJSON(data []byte) error {
	doc := &fieldValueDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	fv.text = doc.Text
	fv.files = doc.Files
	return nil
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
