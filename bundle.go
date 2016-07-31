package fbmodel

import (
	"encoding/json"
	"errors"
)

type Bundle struct {
	ID
	Rev         *string
	Owner       *User
	Name        *string
	Description *string
}

type bundleDoc struct {
	Type        string  `json:"type"`
	ID          ID      `json:"_id"`
	Rev         *string `json:"_rev,omitempty"`
	Owner       string  `json:"owner"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func CreateBundle(id string, owner *User) *Bundle {
	b := &Bundle{}
	b.ID = CreateID("bundle", []byte(id))
	return b
}

func NewBundle(id []byte, owner *User) *Bundle {
	b := &Bundle{}
	b.ID = NewID("bundle", id)
	b.Owner = owner
	return b
}

func (b *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(bundleDoc{
		Type:        "bundle",
		ID:          b.ID,
		Rev:         b.Rev,
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
		return errors.New("Invalid document type for bundle")
	}
	b.Rev = doc.Rev
	user, err := NewUserStub(doc.Owner)
	if err != nil {
		return err
	}
	b.ID = doc.ID
	b.Owner = user
	b.Name = doc.Name
	b.Description = doc.Description

	return nil
}
