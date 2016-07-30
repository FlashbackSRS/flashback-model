package fbmodel

import (
	"errors"
	"strings"
)

var validTypes map[string]struct{}

func init() {
	validTypes = make(map[string]struct{})
	for _, t := range []string{"theme", "note", "deck", "bundle", "card", "user"} {
		validTypes[t] = struct{}{}
	}
}

func isValidType(t string) bool {
	_, ok := validTypes[t]
	return ok
}

type typeID struct {
	docType string
	docID   string
}

func (t *typeID) parse(docID string) error {
	parts := strings.SplitN(docID, "-", 2)
	return t.set(parts[0], parts[1])
}

func (t *typeID) set(tp, id string) error {
	if t.docType != "" || t.docID != "" {
		return errors.New("typeID already set")
	}
	if !isValidType(tp) {
		return errors.New("Invalid document type: " + tp)
	}
	t.docType = tp
	t.docID = id
	return nil
}

func (t *typeID) isValidID(id string) bool {
	return true
}

func (t *typeID) DocID() string {
	return t.docType + "-" + t.docID
}

func (t *typeID) ID() string {
	return t.docID
}

func (t *typeID) Type() string {
	return t.docType
}

type baseDoc struct {
	Type string  `json:"type"`
	ID   string  `json:"_id"`
	Rev  *string `json:"_rev,omitempty"`
}

type namedDoc struct {
	baseDoc
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
