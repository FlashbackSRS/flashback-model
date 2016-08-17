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
	ID
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	Owner       *User
	Name        *string
	Description *string
}

type bundleDoc struct {
	Type        string     `json:"type"`
	ID          ID         `json:"_id"`
	Rev         *string    `json:"_rev,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	Modified    *time.Time `json:"modified,omitempty"`
	Imported    *time.Time `json:"imported,omitempty"`
	Owner       string     `json:"owner"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func NewBundle(id string, owner *User) (*Bundle, error) {
	b := &Bundle{}
	bid, err := NewID("bundle", id)
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
		Owner:       b.Owner.Identity(),
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
