package fbmodel

import ()

type Model struct {
	NamedDoc
	Files *FileCollectionView `json:"files,omitempty"`
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}

func (m *Model) SetFile(name, ctype string, content []byte) {
	m.Files.SetFile(name, ctype, content)
}
