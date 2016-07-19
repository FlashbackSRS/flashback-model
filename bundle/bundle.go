// A bundle is an actual CouchDB database, which contains themes, decks, and notes.
// A bundle is the smallest, sharable unit in Flashback. As such, ownership and
// permissions are defind per bundle
package bundle

import (
	"encoding/json"

	"github.com/flimzy/flashback-model"
	"github.com/flimzy/flashback-model/user"
)

type Bundle struct {
	model.NamedDoc
	Owner *user.User
}

func New(id string, owner *user.User) *Bundle {
	b := &Bundle{}
	b.Type = "bundle"
	b.ID = "bundle-" + id
	b.Owner = owner
	return b
}

type jsonBundle struct {
	model.NamedDoc
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
