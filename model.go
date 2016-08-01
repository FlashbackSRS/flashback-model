package fbmodel

import (
	"strconv"
	// 	"encoding/json"
)

type Model struct {
	Theme       *Theme              `json:"-"`
	ID          uint32              `json:"id"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Files       *FileCollectionView `json:"files,omitempty"`
}

func NewModel(t *Theme) (*Model, error) {
	return &Model{
		Theme: t,
		ID:    t.NextModelSequence(),
		Files: t.Attachments.NewView(),
	}, nil
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}

func (m *Model) Identity() string {
	return m.Theme.ID.Identity() + "." + strconv.FormatUint(uint64(m.ID), 16)
}
