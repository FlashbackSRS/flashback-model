package fbmodel

import (
	"encoding/json"
)

type Bundle struct {
	NamedDoc
	Owner *User
}

func NewBundle(id string, owner *User) *Bundle {
	b := &Bundle{}
	b.Type = "bundle"
	b.ID = "bundle-" + id
	b.Owner = owner
	return b
}

type jsonBundle struct {
	NamedDoc
	Owner string `json:"owner"`
}

// MarshalJSON satisfies the json.Marshaler interface for Bundle structs
func (b *Bundle) MarshalJSON() ([]byte, error) {
	if b.Type != "bundle" {
		panic("Invalid bundle type")
	}
	return json.Marshal(jsonBundle{
		NamedDoc: b.NamedDoc,
		Owner:    b.Owner.UUID().String(),
	})
}
