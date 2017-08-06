package fb

import (
	"bytes"
	"testing"

	"github.com/flimzy/diff"
)

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
