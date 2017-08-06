package fb

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// Theme contains related visual representation elements.
//
// A theme should generally group all card types (models) of a standard visual
// theme. A theme may contain files that are shared across all included modules,
// such as a CSS theme, common graphics files, or even JavaScript. Within a
// theme must exist one or more models, each of which represents a specific card
// type, and which may additionally have its own attachments.
type Theme struct {
	ID            DocID
	Rev           *string
	Created       time.Time
	Modified      time.Time
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
	ID            DocID               `json:"_id"`
	Rev           *string             `json:"_rev,omitempty"`
	Created       time.Time           `json:"created"`
	Modified      time.Time           `json:"modified"`
	Imported      *time.Time          `json:"imported,omitempty"`
	Name          *string             `json:"name,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Models        []*Model            `json:"models,omitempty"`
	Attachments   *FileCollection     `json:"_attachments,omitempty"`
	Files         *FileCollectionView `json:"files,omitempty"`
	ModelSequence uint32              `json:"modelSequence"`
}

// Validate validates that all of the data in the theme appears valid and self
// consistent. A nil return value means no errors were detected.
func (t *themeDoc) Validate() error {
	if t.ID.id == nil || len(t.ID.id) == 0 {
		return errors.New("id required")
	}
	if t.ID.docType != "theme" {
		return errors.New("invalid doc type")
	}
	if t.Created.IsZero() {
		return errors.New("created time required")
	}
	if t.Modified.IsZero() {
		return errors.New("modified time required")
	}
	if t.Attachments == nil {
		return errors.New("attachments collection must not be nil")
	}
	if t.Files == nil {
		return errors.New("file list must not be nil")
	}
	if !t.Attachments.hasMemberView(t.Files) {
		return errors.New("file list must be a member of attachments collection")
	}
	for _, m := range t.Models {
		if m.ID <= t.ModelSequence {
			return errors.New("modelSequence must larger than existing model IDs")
		}
	}
	return nil
}

// NewTheme returns a new, bare-bones theme, with the specified ID.
func NewTheme(id []byte) (*Theme, error) {
	t := &Theme{}
	tid, err := NewDocID("theme", id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DocID for Theme")
	}
	t.ID = tid
	t.Attachments = NewFileCollection()
	t.Files = t.Attachments.NewView()
	t.Models = make([]*Model, 0, 1)
	return t, nil
}

// SetFile sets an attachment with the requested name, type, and content, as
// part of the Theme, overwriting any attachment with the same name, if it exists.
func (t *Theme) SetFile(name, ctype string, content []byte) {
	t.Files.SetFile(name, ctype, content)
}

// MarshalJSON implements the json.Marshaler interface for the Theme type.
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

// NewModel returns a new model of the requested type.
func (t *Theme) NewModel(mType string) (*Model, error) {
	m, err := NewModel(t, mType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create model")
	}
	t.Models = append(t.Models, m)
	return m, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Theme type.
func (t *Theme) UnmarshalJSON(data []byte) error {
	doc := &themeDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal Theme")
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

	if t.Attachments == nil {
		return errors.New("invalid theme: no attachments")
	}
	if t.Files == nil {
		return errors.New("invalid theme: no file list")
	}

	if err := t.Attachments.AddView(t.Files); err != nil {
		return err
	}
	for _, m := range t.Models {
		if err := t.Attachments.AddView(m.Files); err != nil {
			return err
		}
		m.Theme = t
	}

	return nil
}

// NextModelSequence returns the next available model sequence, while also
// updating the internal counter.
func (t *Theme) NextModelSequence() uint32 {
	id := t.modelSequence
	atomic.AddUint32(&t.modelSequence, 1)
	return id
}

// SetRev sets the _rev attribute of the Theme.
func (t *Theme) SetRev(rev string) { t.Rev = &rev }

// DocID returns the theme's _id
func (t *Theme) DocID() string { return t.ID.String() }

// ImportedTime returns the time the Theme was imported, or nil
func (t *Theme) ImportedTime() *time.Time { return t.Imported }

// ModifiedTime returns the time the Theme was last modified
func (t *Theme) ModifiedTime() *time.Time { return &t.Modified }

// MergeImport attempts to merge i into t and returns true if a merge occurred,
// or false if no merge was necessary.
func (t *Theme) MergeImport(i interface{}) (bool, error) {
	existing := i.(*Theme)
	if !t.ID.Equal(&existing.ID) {
		return false, errors.New("IDs don't match")
	}
	if t.Imported == nil || existing.Imported == nil {
		return false, errors.New("not an import")
	}
	if !t.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	t.Rev = existing.Rev
	if t.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	t.Name = existing.Name
	t.Description = existing.Description
	t.Models = existing.Models
	t.Attachments = existing.Attachments
	t.Files = existing.Files
	t.modelSequence = existing.modelSequence
	t.Modified = existing.Modified
	t.Imported = existing.Imported
	return false, nil
}
