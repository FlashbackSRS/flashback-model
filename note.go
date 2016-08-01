package fbmodel

import (
	"encoding/json"
	"time"
)

type Note struct {
	ID
	Rev      *string
	Created  *time.Time
	Modified *time.Time
	Imported *time.Time
	Model    *Model
}

type noteDoc struct {
	Type     string     `json:"type"`
	ID       ID         `json:"_id"`
	Rev      *string    `json:"_rev,omitempty"`
	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Imported *time.Time `json:"imported,omitempty"`
	Model    string     `json:"model"`
}

func NewNote(id string, model *Model) (*Note, error) {
	n := &Note{}
	nid, err := NewID("note", id)
	if err != nil {
		return nil, err
	}
	n.ID = nid
	n.Model = model
	return n, nil
}

func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(noteDoc{
		Type:     "note",
		ID:       n.ID,
		Rev:      n.Rev,
		Created:  n.Created,
		Modified: n.Modified,
		Imported: n.Imported,
		Model:    n.Model.Identity(),
	})
}

/*
type Note struct {
    Tags           string           `db:"tags"` // List of the note's tags
    FieldValues    FieldValues      `db:"flds"` // Values for the note's fields
    UniqueField    string           `db:"sfld"` // The text of the first field, used for Anki's simplistic uniqueness checking
    Checksum       int64            `db:"csum"` // Field checksum used for duplicate check. Integer representation of first 8 digits of sha1 hash of the first field
}
*/
