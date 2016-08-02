package fb

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"
)

type Theme struct {
	ID
	Rev           *string
	Created       *time.Time
	Modified      *time.Time
	Imported      *time.Time
	Name          *string
	Description   *string
	Models        []*Model
	Attachments   *FileCollection
	Files         *FileCollectionView
	modelSequence uint32
}

type themeDoc struct {
	Type          string              `json:"type"`
	ID            ID                  `json:"_id"`
	Rev           *string             `json:"_rev,omitempty"`
	Created       *time.Time          `json:"created,omitempty"`
	Modified      *time.Time          `json:"modified,omitempty"`
	Imported      *time.Time          `json:"imported,omitempty"`
	Name          *string             `json:"name,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Models        []*Model            `json:"models,omitempty"`
	Attachments   *FileCollection     `json:"_attachments,omitempty"`
	Files         *FileCollectionView `json:"files,omitempty"`
	ModelSequence uint32              `json:"modelSequence"`
}

func NewTheme(id string) (*Theme, error) {
	t := &Theme{}
	tid, err := NewID("theme", id)
	if err != nil {
		return nil, err
	}
	t.ID = tid
	t.Attachments = NewFileCollection()
	t.Files = t.Attachments.NewView()
	t.Models = make([]*Model, 0, 1)
	return t, nil
}

func (t *Theme) SetFile(name, ctype string, content []byte) {
	t.Files.SetFile(name, ctype, content)
}

func (t *Theme) MarshalJSON() ([]byte, error) {
	return json.Marshal(themeDoc{
		Type:          "theme",
		ID:            t.ID,
		Rev:           t.Rev,
		Created:       t.Created,
		Modified:      t.Modified,
		Imported:      t.Imported,
		Name:          t.Name,
		Description:   t.Description,
		Models:        t.Models,
		Attachments:   t.Attachments,
		Files:         t.Files,
		ModelSequence: t.modelSequence,
	})
}

func (t *Theme) NewModel(mType ModelType) (*Model, error) {
	m, err := NewModel(t, mType)
	if err != nil {
		return nil, err
	}
	t.Models = append(t.Models, m)
	return m, nil
}

func (t *Theme) UnmarshalJSON(data []byte) error {
	doc := &themeDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "theme" {
		return errors.New("Invalid document type for theme: " + doc.Type)
	}
	t.ID = doc.ID
	t.Rev = doc.Rev
	t.Created = doc.Created
	t.Modified = doc.Modified
	t.Imported = doc.Imported
	t.Name = doc.Name
	t.Description = doc.Description
	t.Models = doc.Models
	t.Attachments = doc.Attachments
	t.Files = doc.Files
	t.modelSequence = doc.ModelSequence

	t.Attachments.AddView(t.Files)
	for _, m := range t.Models {
		t.Attachments.AddView(m.Files)
		m.Theme = t
	}

	return nil
}

func (t *Theme) NextModelSequence() uint32 {
	id := t.modelSequence
	atomic.AddUint32(&t.modelSequence, 1)
	return id
}
