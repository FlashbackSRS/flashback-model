package fbmodel

import (
	"encoding/json"
	"time"
)

type Theme struct {
	ID
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	Name        *string
	Description *string
	Models      []*Model
	Attachments *FileCollection
	Files       *FileCollectionView
}

type themeDoc struct {
	Type        string              `json:"type"`
	ID          ID                  `json:"_id"`
	Rev         *string             `json:"_rev,omitempty"`
	Created     *time.Time          `json:"created,omitempty"`
	Modified    *time.Time          `json:"modified,omitempty"`
	Imported    *time.Time          `json:"imported,omitempty"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Models      []*Model            `json:"models,omitempty"`
	Attachments *FileCollection     `json:"_attachments,omitempty"`
	Files       *FileCollectionView `json:"files,omitempty"`
}

func CreateTheme(key []byte) *Theme {
	t := &Theme{}
	t.ID = CreateID("theme", key)
	t.Attachments = NewFileCollection()
	t.Files = t.Attachments.NewView()
	t.Models = make([]*Model, 0, 1)
	return t
}

func (t *Theme) SetFile(name, ctype string, content []byte) {
	t.Files.SetFile(name, ctype, content)
}

func (t *Theme) MarshalJSON() ([]byte, error) {
	return json.Marshal(themeDoc{
		Type:        "theme",
		ID:          t.ID,
		Rev:         t.Rev,
		Created:     t.Created,
		Modified:    t.Modified,
		Imported:    t.Imported,
		Name:        t.Name,
		Description: t.Description,
		Models:      t.Models,
		Attachments: t.Attachments,
		Files:       t.Files,
	})
}

func (t *Theme) NewModel(id string) (*Model, error) {
	m, err := NewModel(id, t)
	if err != nil {
		return nil, err
	}
	t.Models = append(t.Models, m)
	return m, nil
}
