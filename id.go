package fbmodel

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

type ID struct {
	docType string
	id      []byte
}

func CreateID(idType string, key []byte) ID {
	hash := sha1.Sum(key)
	return NewID(idType, hash[:])
}

func NewID(idType string, objectID []byte) ID {
	if !isValidType(idType) {
		panic("Invalid type: " + idType)
	}
	id := ID{
		docType: idType,
		id:      objectID,
	}
	return id
}

func (id *ID) String() string {
	return id.docType + "-" + id.Identity()
}

func ParseID(identity string) (*ID, error) {
	parts := strings.SplitN(identity, "-", 2)
	if !isValidType(parts[0]) {
		return nil, errors.New("Invalid type: " + parts[0])
	}
	data, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("Cannot decode ID: " + err.Error())
	}
	return &ID{
		docType: parts[0],
		id:      data,
	}, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	newID, err := ParseID(strings.Trim(string(data), "\""))
	id.docType = newID.docType
	id.id = newID.id
	return err
}

func (id *ID) Identity() string {
	return base64.URLEncoding.EncodeToString(id.id)
}

func (id *ID) Type() string {
	return id.docType
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}
