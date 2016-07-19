package theme

import (
	"github.com/flimzy/flashback-model"
)

const HTMLTemplateContentType = "text/html+flashbacktmpl"

type Theme struct {
	model.NamedDoc
	Models      []*Model                  `json:"models,omitempty"`
	Attachments *model.FileCollection     `json:"_attachments,omitempty"`
	Files       *model.FileCollectionView `json:"files,omitempty"`
}

type Model struct {
	model.NamedDoc
	Files *model.FileCollectionView `json:"files,omitempty"`
}

func NewTheme(id string) *Theme {
	t := &Theme{}
	t.Type = "theme"
	t.ID = "theme-" + id
	t.Models = make([]*Model, 0, 1)
	t.Attachments = model.NewFileCollection()
	t.Files = t.Attachments.NewView()
	return t
}

func (t *Theme) AddFile(name, ctype string, content []byte) error {
	return t.Files.AddFile(name, ctype, content)
}

func (t *Theme) SetFile(name, ctype string, content []byte) {
	t.Files.SetFile(name, ctype, content)
}

func (t *Theme) NewModel(id string) *Model {
	m := &Model{}
	m.Type = "model"
	m.ID = id
	m.Files = t.Attachments.NewView()
	t.Models = append(t.Models, m)
	return m
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}

func (m *Model) SetFile(name, ctype string, content []byte) {
	m.Files.SetFile(name, ctype, content)
}
