package fbmodel

import ()

type Note struct {
	doc
	ModelID string `json:"modelId"`
}

func NewNote(id, modelID string) *Note {
	n := &Note{}
	n.doc = NewDoc("note", id)
	n.ModelID = modelID
	return n
}

/*
type Note struct {
    ModelID        ID               `db:"mid"`  // Model ID
    Modified       TimestampSeconds `db:"mod"`  // Last modified time
    UpdateSequence int              `db:"usn"`  // Update sequence number (no longer used?)
    Tags           string           `db:"tags"` // List of the note's tags
    FieldValues    FieldValues      `db:"flds"` // Values for the note's fields
    UniqueField    string           `db:"sfld"` // The text of the first field, used for Anki's simplistic uniqueness checking
    Checksum       int64            `db:"csum"` // Field checksum used for duplicate check. Integer representation of first 8 digits of sha1 hash of the first field
}
*/
