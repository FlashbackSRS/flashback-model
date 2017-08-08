package fb

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
)

func TestValidateDocID(t *testing.T) {
	type Test struct {
		name string
		id   string
		err  string
	}
	tests := []Test{
		{
			name: "bogus id",
			id:   "really really bogus",
			err:  "invalid DocID format",
		},
		{
			name: "unsupported type",
			id:   "foo-chicken",
			err:  "unsupported DocID type 'foo'",
		},
		{
			name: "invalid base64",
			id:   "deck- really bad stuff",
			err:  "invalid DocID encoding",
		},
		{
			name: "valid",
			id:   "deck-0123456789",
		},
		{
			name: "multiple dashes",
			id:   "deck--v4-v4",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateDocID(test.id)
			checkErr(t, test.err, err)
		})
	}
}

func TestNewDocID(t *testing.T) {
	type Test struct {
		name     string
		doctype  string
		id       []byte
		expected DocID
		err      string
	}
	tests := []Test{
		{
			name:    "invalid type",
			doctype: "chicken",
			err:     "invalid document type: chicken",
		},
		{
			name:    "no id",
			doctype: "deck",
			err:     "id is required",
		},
		{
			name:     "valid",
			doctype:  "deck",
			id:       []byte("test id"),
			expected: DocID{docType: "deck", id: []byte("test id")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewDocID(test.doctype, test.id)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestID(t *testing.T) {
	id, err := NewDocID("note", []byte("Test Note"))
	if err != nil {
		t.Fatalf("Error creating note ID: %s", err)
	}
	if id.docType != "note" {
		t.Fatalf("Unexpected doctype in note id: %s", id.docType)
	}
	if !bytes.Equal([]byte("Test Note"), id.id) {
		t.Fatalf("note ID not as expected")
	}
}

func TestDbID(t *testing.T) {
	id, err := NewDbID("user", []byte("Test User"))
	if err != nil {
		t.Fatalf("Error creating user ID: %s", err)
	}
	if id.docType != "user" {
		t.Fatalf("Unexpected doctype in user id: %s", id.docType)
	}
	if !bytes.Equal([]byte("Test User"), id.id) {
		t.Fatalf("user ID not as expected")
	}
}

func TestEncodeDocID(t *testing.T) {
	expected := "foo-dGVzdCBpZA"
	result := EncodeDocID("foo", []byte("test id"))
	if result != expected {
		t.Errorf("Unexpected result: %s", result)
	}
}
