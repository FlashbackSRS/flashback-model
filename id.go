package fb

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

func KeyToIDString(idType IDType, key ...[]byte) string {
	id := KeyToID(key...)
	return encodeID(id, idType)
}

func encodeID(id []byte, idType IDType) string {
	switch idType {
	case Base64ID:
		return "0" + b64encoder.EncodeToString(id)
	case HexID:
		return "1" + hex.EncodeToString(id)
	default:
		panic(fmt.Sprintf("Invalid ID type: %d\n", idType))
	}
	return ""
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
	idType  IDType
}

type IDType int

const (
	Base64ID IDType = iota
	HexID
)

func isValidIDType(idType IDType) bool {
	return idType >= Base64ID && idType <= HexID
}

// CreateID creates a new ID object, based on the type and a byte array, which is hashed to generate the human-readable ID.
func CreateID(docType string, key []byte) (ID, error) {
	return NewByteID(docType, KeyToID(key), Base64ID)
}

// NewID creates a new ID object, based on the type and human-readable id.
func NewID(docType string, objectID string) (ID, error) {
	idType, id, err := DecodeID(objectID)
	if err != nil {
		return ID{}, err
	}
	return NewByteID(docType, id, idType)
}

// NewByteID creates a new ID object, based on the type and id in []byte format
func NewByteID(docType string, id []byte, idType IDType) (ID, error) {
	if !isValidType(docType) {
		return ID{}, fmt.Errorf("Invalid document type `%s`", docType)
	}
	if !isValidIDType(idType) {
		return ID{}, fmt.Errorf("Invalid ID type: %d", idType)
	}
	return ID{
		docType: docType,
		id:      id,
		idType:  idType,
	}, nil
}

func (id *ID) String() string {
	return id.docType + "-" + id.Identity()
}

func DecodeID(identity string) (IDType, []byte, error) {
	idType := IDType(identity[0] - 48)
	var data []byte
	var err error
	switch idType {
	case Base64ID:
		data, err = b64encoder.DecodeString(identity[1:])
		if err != nil {
			fmt.Printf("\tDecoded: %s, %s\n", identity[1:], err)
		}
	case HexID:
		data, err = hex.DecodeString(identity[1:])
	}
	return idType, data, err
}

func ParseID(identity string) (*ID, error) {
	parts := strings.SplitN(identity, "-", 2)
	idType, data, err := DecodeID(parts[1])
	if err != nil {
		return nil, errors.New("Cannot decode ID: " + err.Error())
	}
	return &ID{
		docType: parts[0],
		id:      data,
		idType:  idType,
	}, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	newID, err := ParseID(strings.Trim(string(data), "\""))
	id.docType = newID.docType
	id.id = newID.id
	id.idType = newID.idType
	return err
}

func (id *ID) Identity() string {
	return encodeID(id.id, id.idType)
}

func (id *ID) RawID() []byte {
	return id.id
}

func (id *ID) Type() string {
	return id.docType
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID) Equal(id2 *ID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}
