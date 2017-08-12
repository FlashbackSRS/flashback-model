package fb

import (
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var validDBIDTypes = map[string]struct{}{
	"user":   {},
	"bundle": {},
}

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
	return b32encoding.DecodeString(s)
}

func validateDBID(id string) error {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return errors.New("invalid DBID format")
	}
	if _, ok := validDBIDTypes[parts[0]]; !ok {
		return errors.Errorf("unsupported DBID type '%s'", parts[0])
	}
	if _, err := b32dec(parts[1]); err != nil {
		return errors.New("invalid DBID encoding")
	}
	return nil
}

// EncodeDBID generates a DBID by encoding the docType and Base32-encoding
// the ID. No validation is done of the docType.
func EncodeDBID(docType string, id []byte) string {
	return fmt.Sprintf("%s-%s", docType, b32enc(id))
}

// DBIDToBytes decodes the DBID into its underlying byte representation.
func DBIDToBytes(id string) ([]byte, error) {
	if err := validateDBID(id); err != nil {
		return nil, err
	}
	parts := strings.SplitN(id, "-", 2)
	return b32dec(parts[1])
}
