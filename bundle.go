package fb

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	TemplateContentType = "text/html"
	BundleContentType   = "application/json"
)

type Bundle struct {
	ID          HexID
	Rev         *string
	Created     time.Time
	Modified    time.Time
	Imported    *time.Time
	Owner       *User
	Name        *string
	Description *string
}

type bundleDoc struct {
	Type        string     `json:"type"`
	ID          HexID      `json:"_id"`
	Rev         *string    `json:"_rev,omitempty"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	Imported    *time.Time `json:"imported,omitempty"`
	Owner       string     `json:"owner"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func NewBundle(id []byte, owner *User) (*Bundle, error) {
	b := &Bundle{}
	bid, err := NewHexID("bundle", id)
	if err != nil {
		return nil, err
	}
	b.ID = bid
	b.Owner = owner
	return b, nil
}

func (b *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(bundleDoc{
		Type:        "bundle",
		ID:          b.ID,
		Rev:         b.Rev,
		Created:     b.Created,
		Modified:    b.Modified,
		Imported:    b.Imported,
		Owner:       b.Owner.ID.Identity(),
		Name:        b.Name,
		Description: b.Description,
	})
}

func (b *Bundle) UnmarshalJSON(data []byte) error {
	doc := &bundleDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "bundle" {
		return errors.New("Invalid document type for bundle: " + doc.Type)
	}
	user, err := NewUserStub(doc.Owner)
	if err != nil {
		return err
	}
	b.ID = doc.ID
	b.Rev = doc.Rev
	b.Created = doc.Created
	b.Modified = doc.Modified
	b.Imported = doc.Imported
	b.Owner = user
	b.Name = doc.Name
	b.Description = doc.Description

	return nil
}

// FlashbackDoc interface
func (b *Bundle) SetRev(rev string)        { b.Rev = &rev }
func (b *Bundle) DocID() string            { return b.ID.String() }
func (b *Bundle) ImportedTime() *time.Time { return b.Imported }
func (b *Bundle) ModifiedTime() *time.Time { return &b.Modified }

func (b *Bundle) MergeImport(i interface{}) (bool, error) {
	existing := i.(*Bundle)
	if !b.ID.Equal(&existing.ID) {
		return false, errors.New("IDs don't match")
	}
	if !b.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	if !b.Owner.Equal(existing.Owner.uuid) {
		return false, errors.New("Cannot change bundle ownership")
	}
	b.Rev = existing.Rev
	if b.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	b.Name = existing.Name
	b.Description = existing.Description
	return false, nil
}
