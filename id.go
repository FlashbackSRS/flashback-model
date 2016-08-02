package fb

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func KeyToID(key ...[]byte) []byte {
	h := sha1.New()
	for _, k := range key {
		h.Write(k)
	}
	hash := h.Sum(nil)
	return hash[:]
}

func KeyToIDString(key ...[]byte) string {
	return base64.URLEncoding.EncodeToString(KeyToID(key...))
}

var validTypes map[string]struct{}

func init() {
	validTypes = make(map[string]struct{})
	for _, t := range []string{"theme", "model", "note", "deck", "bundle", "card", "user"} {
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

// CreateID creates a new ID object, based on the type and a byte array, which is hashed to generate the human-readable ID.
func CreateID(idType string, key []byte) ID {
	return NewByteID(idType, KeyToID(key))
}

// NewID creates a new ID object, based on the type and human-readable id.
func NewID(idType string, objectID string) (ID, error) {
	dec, err := base64.URLEncoding.DecodeString(objectID)
	if err != nil {
		return ID{}, err
	}
	return NewByteID(idType, dec), nil
}

// NewByteID creates a new ID object, based on the type and id in []byte format
func NewByteID(idType string, id []byte) ID {
	if !isValidType(idType) {
		panic("Invalid type: " + idType)
	}
	return ID{
		docType: idType,
		id:      id,
	}
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
