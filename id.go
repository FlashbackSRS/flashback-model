package fb

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
)

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

type baseID struct {
	docType string
	id      []byte
}

// A standard Base64-encoded ID
type ID struct {
	baseID
}

func newBaseID(docType string, id []byte) (baseID, error) {
	if !isValidType(docType) {
		return baseID{}, errors.New("Invalid document type:" + docType)
	}
	return baseID{
		docType: docType,
		id:      id,
	}, nil
}

func (id *baseID) Type() string {
	return id.docType
}

func parseParts(input ...string) (string, string) {
	switch len(input) {
	case 1:
		parts := strings.SplitN(input[0], "-", 2)
		return parts[0], parts[1]
	case 2:
		return input[0], input[1]
	default:
		panic("IDs must have exactly 1 or 2 parts")
	}
	return "", ""
}

func ParseID(parts ...string) (ID, error) {
	id := ID{}
	err := id.parse(parts...)
	return id, err
}

func (id *ID) parse(parts ...string) error {
	docType, identity := parseParts(parts...)
	data, err := b64encoder.DecodeString(identity)
	if err != nil {
		return err
	}
	base, err := newBaseID(docType, data)
	if err != nil {
		return err
	}
	id.baseID = base
	return nil
}

func NewID(docType string, id []byte) (ID, error) {
	base, err := newBaseID(docType, id)
	return ID{base}, err
}

func (id ID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

func (id *ID) String() string {
	return id.docType + "-" + id.Identity()
}

func (id *ID) Identity() string {
	return b64encoder.EncodeToString(id.id)
}

func (id *ID) Equal(id2 *ID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

// A Hex-encoded ID, for documents which need their own databases (which don't support Base64 alphabets)
type HexID struct {
	baseID
}

func ParseHexID(parts ...string) (HexID, error) {
	docType, identity := parseParts(parts...)
	data, err := hex.DecodeString(identity)
	if err != nil {
		return HexID{}, err
	}
	id, err := newBaseID(docType, data)
	return HexID{id}, err
}

func (id *HexID) parse(parts ...string) error {
	docType, identity := parseParts(parts...)
	data, err := hex.DecodeString(identity)
	if err != nil {
		return err
	}
	base, err := newBaseID(docType, data)
	if err != nil {
		return err
	}
	id.baseID = base
	return nil
}

func NewHexID(docType string, id []byte) (HexID, error) {
	base, err := newBaseID(docType, id)
	return HexID{base}, err
}

func (id HexID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

func (id *HexID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

func (id *HexID) String() string {
	return id.docType + "-" + id.Identity()
}

func (id *HexID) Identity() string {
	return hex.EncodeToString(id.id)
}

func (id *HexID) Equal(id2 *HexID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

func (id *HexID) RawID() []byte {
	return id.id
}
