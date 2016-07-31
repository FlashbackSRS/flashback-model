package fbmodel

import (
	"encoding/json"
	"time"
)

type Model struct {
	ID
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	Name        *string
	Description *string
	Files       *FileCollectionView
}

type modelDoc struct {
	Type        string              `json:"type"`
	ID          ID                  `json:"_id"`
	Rev         *string             `json:"_rev,omitempty"`
	Created     *time.Time          `json:"created,omitempty"`
	Modified    *time.Time          `json:"modified,omitempty"`
	Imported    *time.Time          `json:"imported,omitempty"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Files       *FileCollectionView `json:"files,omitempty"`
}

func NewModel(id string, t *Theme) (*Model, error) {
	m := &Model{}
	if id, err := NewID("model", id); err != nil {
		return nil, err
	} else {
		m.ID = id
	}
	m.Files = t.Attachments.NewView()
	return m, nil
}

func (m *Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(modelDoc{
		Type:        "model",
		ID:          m.ID,
		Rev:         m.Rev,
		Created:     m.Created,
		Modified:    m.Modified,
		Imported:    m.Imported,
		Name:        m.Name,
		Description: m.Description,
	})
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}
