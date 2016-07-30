package fbmodel

import (
	"encoding/json"
)

type Bundle struct {
	namedDoc
	Owner *User
}

func NewBundle(id string, owner *User) *Bundle {
	b := &Bundle{}
	b.doc = NewDoc("bundle", id)
	b.Owner = owner
	return b
}

type jsonBundle struct {
	*jsonDoc
	*Bundle
}

func (b *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonBundle{
		jsonDoc: b.jsonDoc(),
		Bundle: b,
	})
}
