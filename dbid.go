package fb

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var validDbIDTypes map[string]struct{}

// Same as standard Base32 encoding, only lowercase to work with CouchDB database
// naming restrictions.
var b32encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

func b32enc(data []byte) string {
	return strings.TrimRight(b32encoding.EncodeToString(data), "=")
}

func b32dec(s string) ([]byte, error) {
	// fmt.Printf("Before: '%s'\n", s)
	if padLen := len(s) % 8; padLen > 0 {
		s = s + strings.Repeat("=", 8-padLen)
	}
	// fmt.Printf(" After: '%s'\n", s)
	data, err := b32encoding.DecodeString(s)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	return data, err
}

func init() {
	validDbIDTypes = make(map[string]struct{})
	for _, t := range []string{"bundle", "user"} {
		validDbIDTypes[t] = struct{}{}
	}
}

func isValidDbIDType(t string) bool {
	_, ok := validDbIDTypes[t]
	return ok
}

// DbID represents a standard document ID
// Valid characters: a-z, 0-9, _, $, (, ), +, - (But I prefer to avoid -)
// See http://wiki.apache.org/couchdb/HTTP_database_API#Naming_and_Addressing
type DbID struct {
	docType string
	id      []byte
}

// Type returns the DbID's docType.
func (id *DbID) Type() string {
	return id.docType
}

// ParseDbID parses a string representation of a DbID, returning the DbID.
func ParseDbID(parts ...string) (DbID, error) {
	docType, identity := parseParts(parts...)
	data, err := b32dec(identity)
	if err != nil {
		return DbID{}, err
	}
	return NewDbID(docType, data)
}

func (id *DbID) parse(parts ...string) error {
	docType, identity := parseParts(parts...)
	data, err := b32dec(identity)
	if err != nil {
		return errors.Wrap(err, "invalid Base32 in DbID")
	}
	if !isValidDbIDType(docType) {
		return errors.New("Invalid docType: " + docType)
	}
	id.docType = docType
	id.id = data
	return nil
}

// NewDbID returns a new DbID based on the specified docType and id.
func NewDbID(docType string, id []byte) (DbID, error) {
	if !isValidDbIDType(docType) {
		return DbID{}, errors.New("Invalid document type:" + docType)
	}
	return DbID{
		docType: docType,
		id:      id,
	}, nil
}

// MarshalJSON implements the json.Marshaler interface for the DbID type.
func (id DbID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for the DbID type.
func (id *DbID) UnmarshalJSON(data []byte) error {
	raw := string(data)
	return id.parse(raw[1 : len(raw)-1])
}

// String returns the DbID's full string representation.
func (id *DbID) String() string {
	return id.docType + "-" + id.Identity()
}

// Identity returns the DbID's identity as a string.
func (id *DbID) Identity() string {
	return b32enc(id.id)
}

// Equal returns true if both DbIDs are equal.
func (id *DbID) Equal(id2 *DbID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

// RawID returns the byte representation of the internal identity.
func (id *DbID) RawID() []byte {
	return id.id
}
