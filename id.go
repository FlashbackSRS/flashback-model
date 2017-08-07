package fb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var validDocIDTypes = map[string]struct{}{
	"theme": {},
	"model": {},
	"note":  {},
	"deck":  {},
	"card":  {},
}

func validateDocID(id string) error {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return errors.New("invalid DocID format")
	}
	if _, ok := validDocIDTypes[parts[0]]; !ok {
		return errors.Errorf("unsupported DocID type '%s'", parts[0])
	}
	if _, err := b64encoder.DecodeString(parts[1]); err != nil {
		return errors.New("invalid DocID encoding")
	}
	return nil
}

func isValidDocIDType(t string) bool {
	_, ok := validDocIDTypes[t]
	return ok
}

// DocID represents a standard document ID
// Current implementation uses Base64
type DocID struct {
	docType string
	id      []byte
}

// Type returns the DocID's docType.
func (id *DocID) Type() string {
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

// ParseDocID parses a string reprsentation of a DocID, returning the DocID or an error.
func ParseDocID(parts ...string) (DocID, error) {
	id := DocID{}
	err := id.parse(parts...)
	return id, errors.Wrap(err, "failed to parse DocID")
}

func (id *DocID) parse(parts ...string) error {
	docType, identity := parseParts(parts...)
	data, err := b64encoder.DecodeString(identity)
	if err != nil {
		return errors.Wrap(err, "invalid Base64")
	}
	if !isValidDocIDType(docType) {
		return errors.New("Invalid docType: " + docType)
	}
	id.docType = docType
	id.id = data
	return nil
}

// EncodeDocID generates a DocID by encoding the docType and Base64-encoding
// the ID. No validation is done of the docType.
func EncodeDocID(docType string, id []byte) string {
	return fmt.Sprintf("%s-%s", docType, b64encoder.EncodeToString(id))
}

// NewDocID returns a new ID with the provided docType and Identity.
func NewDocID(docType string, id []byte) (DocID, error) {
	if !isValidDocIDType(docType) {
		return DocID{}, errors.New("invalid document type: " + docType)
	}
	if id == nil || len(id) == 0 {
		return DocID{}, errors.New("id is required")
	}
	return DocID{
		docType: docType,
		id:      id,
	}, nil
}

// MarshalJSON implements the json.Marshaler interface for the DocID type.
func (id DocID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the ID type.
func (id *DocID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

// String returns the full string representation of the ID.
func (id *DocID) String() string {
	return id.docType + "-" + id.Identity()
}

// Identity returns a string representation of the internal identity only.
func (id *DocID) Identity() string {
	return b64encoder.EncodeToString(id.id)
}

// Equal returns true if the two DocIDs are considered equal.
func (id *DocID) Equal(id2 *DocID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}
