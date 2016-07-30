package fbmodel

import ()

type typeID struct {
	id      string
	docType string
}

type baseDoc struct {
	ID   string  `json:"_id"`
	Rev  *string `json:"_rev,omitempty"`
	Type string  `json:"type"`
}

type namedDoc struct {
	baseDoc
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type doc struct {
	rev     *string
	id      string
	docType string
}

func (d *doc) ID() string {
	return d.id
}

func (d *doc) Type() string {
	return d.docType
}

func (d *doc) Rev() *string {
	return d.rev
}

func (d *doc) jsonDoc() *jsonDoc {
	return &jsonDoc{
		Rev:  d.Rev(),
		ID:   d.Type() + "-" + d.ID(),
		Type: d.Type(),
	}
}

func (d *doc) MarshalJSON() ([]byte, error) {
	panic("You must provide a MarshalJSON() method for docType " + d.docType)
}

var validTypes map[string]struct{}

func init() {
	validTypes = make(map[string]struct{})
	for _, t := range []string{"theme", "note", "deck", "bundle", "card", "user"} {
		validTypes[t] = struct{}{}
	}
}

func NewDoc(docType, id string) doc {
	if _, ok := validTypes[docType]; !ok {
		panic(docType + " is not a valid docType!")
	}
	return doc{
		id:      id,
		docType: docType,
	}
}

func (d *doc) SetRev(rev string) {
	d.rev = &rev
}
