package fb

import (
	"strconv"
	// 	"encoding/json"
)

type Model struct {
	Theme       *Theme              `json:"-"`
	ID          uint32              `json:"id"`
	Type        ModelType           `json:"modelType"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Fields      []*Field            `json:"fields"`
	Files       *FileCollectionView `json:"files,omitempty"`
}

type ModelType int

const (
	AnkiStandard ModelType = iota
	AnkiCloze
)

func NewModel(t *Theme, mType ModelType) (*Model, error) {
	return &Model{
		Theme:  t,
		ID:     t.NextModelSequence(),
		Type:   mType,
		Fields: make([]*Field, 0, 1),
		Files:  t.Attachments.NewView(),
	}, nil
}

func (m *Model) AddFile(name, ctype string, content []byte) error {
	return m.Files.AddFile(name, ctype, content)
}

func (m *Model) Identity() string {
	return m.Theme.ID.Identity() + "." + strconv.FormatUint(uint64(m.ID), 16)
}

func (m *Model) AddField(fType FieldType, name string) error {
	m.Fields = append(m.Fields, &Field{
		Type: fType,
		Name: name,
	})
	return nil
}

type FieldType int

const (
	TextField FieldType = iota
	ImageField
	AudioField
	AnkiField
)

// A field of a model
//
// Excluded from this definition is the `media` field, which appears to no longer be used.
type Field struct {
	Type FieldType `json:"fieldType"`
	Name string    `json:"name"` // Field name
}
