package fbmodel

type Theme struct {
	namedDoc
	Models      []*Model            `json:"models,omitempty"`
	Attachments *FileCollection     `json:"_attachments,omitempty"`
	Files       *FileCollectionView `json:"files,omitempty"`
}

func NewTheme(id string) *Theme {
	t := &Theme{}
	t.doc = NewDoc("theme", id)
	t.Models = make([]*Model, 0, 1)
	t.Attachments = NewFileCollection()
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
	m.doc = NewDoc("model", id)
	m.Files = t.Attachments.NewView()
	t.Models = append(t.Models, m)
	return m
}
