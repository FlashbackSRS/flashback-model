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

// ID represents a standard Base64-encoded ID
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
}

// ParseID parses a string reprsentation of an ID, returning the ID or an error.
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

// NewID returns a new ID with the provided docType and Identity.
func NewID(docType string, id []byte) (ID, error) {
	base, err := newBaseID(docType, id)
	return ID{base}, err
}

// MarshalJSON implements the json.Marshaler interface for the ID type.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the ID type.
func (id *ID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

// String returns the full string representation of the ID.
func (id *ID) String() string {
	return id.docType + "-" + id.Identity()
}

// Identity returns a string representation of the internal ID only.
func (id *ID) Identity() string {
	return b64encoder.EncodeToString(id.id)
}

// Equal returns true if the two IDs are considered equal.
func (id *ID) Equal(id2 *ID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

// HexID represents a Hex-encoded ID, for documents which need their own databases
// (which don't support Base64 alphabets)
type HexID struct {
	baseID
}

// ParseHexID parses a string representation of a HexID, returning the HexID.
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

// NewHexID returns a new HexID based on the specified docType and id.
func NewHexID(docType string, id []byte) (HexID, error) {
	base, err := newBaseID(docType, id)
	return HexID{base}, err
}

// MarshalJSON implements the json.Marshaler interface for the HexID type.
func (id HexID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the HexID type.
func (id *HexID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

// String returns the HexID's full string representation.
func (id *HexID) String() string {
	return id.docType + "-" + id.Identity()
}

// Identity returns the HexID's identity as a string.
func (id *HexID) Identity() string {
	return hex.EncodeToString(id.id)
}

// Equal returns true if both HexIDs are equal.
func (id *HexID) Equal(id2 *HexID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

// RawID returns the byte representation of the internal identity.
func (id *HexID) RawID() []byte {
	return id.id
}
