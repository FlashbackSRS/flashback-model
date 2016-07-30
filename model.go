package fbmodel

import ()

type Model struct {
	namedDoc
	Files *FileCollectionView `json:"files,omitempty"`
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}

func (m *Model) SetFile(name, ctype string, content []byte) {
	m.Files.SetFile(name, ctype, content)
}

// func (m *Model) NewNote() *Note {
// 	n := &Note{}
// 	n.ModelID = m.ID
// 	n.Type = "note"
// 	
// }

/*
type Note struct {
    ID             ID               `db:"id"`   // Primary key
    GUID           string           `db:"guid"` // globally unique id, almost certainly used for syncing
    ModelID        ID               `db:"mid"`  // Model ID
    Modified       TimestampSeconds `db:"mod"`  // Last modified time
    UpdateSequence int              `db:"usn"`  // Update sequence number (no longer used?)
    Tags           string           `db:"tags"` // List of the note's tags
    FieldValues    FieldValues      `db:"flds"` // Values for the note's fields
    UniqueField    string           `db:"sfld"` // The text of the first field, used for Anki's simplistic uniqueness checking
    Checksum       int64            `db:"csum"` // Field checksum used for duplicate check. Integer representation of first 8 digits of sha1 hash of the first field
}
*/
