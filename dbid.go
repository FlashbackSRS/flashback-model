package fb

import (
	"bytes"
	"encoding/hex"
	"errors"
)

var validDbIDTypes map[string]struct{}

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
// Valid characters: a-z, 0-9, _, $, (, ), + -
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
	data, err := hex.DecodeString(identity)
	if err != nil {
		return DbID{}, err
	}
	return NewDbID(docType, data)
}

func (id *DbID) parse(parts ...string) error {
	docType, identity := parseParts(parts...)
	data, err := hex.DecodeString(identity)
	if err != nil {
		return err
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
	return hex.EncodeToString(id.id)
}

// Equal returns true if both DbIDs are equal.
func (id *DbID) Equal(id2 *DbID) bool {
	return id.docType == id2.docType && bytes.Equal(id.id, id2.id)
}

// RawID returns the byte representation of the internal identity.
func (id *DbID) RawID() []byte {
	return id.id
}
