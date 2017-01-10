package fb

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Note represents a Flashback note.
type Note struct {
	ID          DocID
	Rev         *string
	Created     time.Time
	Modified    time.Time
	Imported    *time.Time
	ThemeID     string
	ModelID     uint32
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
	ID          DocID           `json:"_id"`
	Rev         *string         `json:"_rev,omitempty"`
	Created     time.Time       `json:"created"`
	Modified    time.Time       `json:"modified"`
	Imported    *time.Time      `json:"imported,omitempty"`
	ThemeID     string          `json:"theme"`
	ModelID     uint32          `json:"model"`
	FieldValues []*FieldValue   `json:"fieldValues"`
	Attachments *FileCollection `json:"_attachments,omitempty"`
}

// NewNote creates a new, empty note with the provided ID and Model.
func NewNote(id []byte, model *Model) (*Note, error) {
	n := &Note{}
	nid, err := NewDocID("note", id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create DocID")
	}
	n.ID = nid
	n.ThemeID = model.Theme.ID.String()
	n.ModelID = model.ID
	n.FieldValues = make([]*FieldValue, len(model.Fields))
	n.Attachments = NewFileCollection()
	n.model = model
	return n, nil
}

// SetModel assigns the provided model to the Note. This is useful after retrieving
// a note.
func (n *Note) SetModel(m *Model) {
	n.model = m
	for i := 0; i < len(n.FieldValues); i++ {
		n.FieldValues[i].field = m.Fields[i]
	}
}

// Model returns the Note's associated Model.
func (n *Note) Model() *Model {
	return n.model
}

// MarshalJSON implements the json.Marshaler interface for the Note type.
func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(noteDoc{
		Type:        "note",
		ID:          n.ID,
		Rev:         n.Rev,
		Created:     n.Created,
		Modified:    n.Modified,
		Imported:    n.Imported,
		ThemeID:     n.ThemeID,
		ModelID:     n.ModelID,
		FieldValues: n.FieldValues,
		Attachments: n.Attachments,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Note type.
func (n *Note) UnmarshalJSON(data []byte) error {
	doc := &noteDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal Note")
	}
	if doc.Type != "note" {
		return errors.New("Invalid document type for note: " + doc.Type)
	}
	n.ID = doc.ID
	n.Rev = doc.Rev
	n.Created = doc.Created
	n.Modified = doc.Modified
	n.Imported = doc.Imported
	n.ThemeID = doc.ThemeID
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

// GetFieldValue returns the requested FieldValue by index.
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

// Type returns the FieldType of the FieldValue.
func (fv *FieldValue) Type() FieldType {
	return fv.field.Type
}

// FieldValue stores the value of a given field.
type FieldValue struct {
	field *Field
	text  *string
	files *FileCollectionView
}

type fieldValueDoc struct {
	Text  *string             `json:"text,omitempty"`
	Files *FileCollectionView `json:"files,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for the FieldValue type.
func (fv *FieldValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(fieldValueDoc{
		Text:  fv.text,
		Files: fv.files,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for the FieldValue type.
func (fv *FieldValue) UnmarshalJSON(data []byte) error {
	doc := &fieldValueDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal FieldValue")
	}
	fv.text = doc.Text
	fv.files = doc.Files
	return nil
}

// SetText sets the text attribute of the FieldValue, or returns an error if the
// FieldValue's type does not permit text attributes.
func (fv *FieldValue) SetText(text string) error {
	if fv.field.Type != TextField && fv.field.Type != AnkiField {
		return errors.New("Text field not permitted")
	}
	fv.text = &text
	return nil
}

// Text retrieves the text of the field value, or an error if no text field
// is permitted for the field type.
func (fv *FieldValue) Text() (string, error) {
	if fv.field.Type != TextField && fv.field.Type != AnkiField {
		return "", errors.New("FieldValue has no text field")
	}
	return *fv.text, nil
}

// AddFile adds a file of the specified name, type, and content, as an attachment
// to be used by the FieldValue.
func (fv *FieldValue) AddFile(name, ctype string, content []byte) error {
	if fv.field.Type == TextField {
		return errors.New("Text fields do not support attachments")
	}
	return fv.files.AddFile(name, ctype, content)
}

// SetRev sets the Note's _rev attribute.
func (n *Note) SetRev(rev string) { n.Rev = &rev }

// DocID returns the Note's _id attribute.
func (n *Note) DocID() string { return n.ID.String() }

// ImportedTime returns the time the Note was imported, or nil.
func (n *Note) ImportedTime() *time.Time { return n.Imported }

// ModifiedTime returns the time the Note was last modified.
func (n *Note) ModifiedTime() *time.Time { return &n.Modified }

// MergeImport attempts to merge i into n, returning true if successful, or
// false if no merge was necessary.
func (n *Note) MergeImport(i interface{}) (bool, error) {
	existing := i.(*Note)
	if !n.ID.Equal(&existing.ID) {
		return false, errors.New("IDs don't match")
	}
	if !n.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	n.Rev = existing.Rev
	if n.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	n.Modified = existing.Modified
	n.Imported = existing.Imported
	n.ModelID = existing.ModelID
	n.FieldValues = existing.FieldValues
	n.Attachments = existing.Attachments
	n.model = existing.model
	return false, nil
}
