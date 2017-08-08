package fb

import (
	"bytes"
	"testing"
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateDocID(test.id)
			checkErr(t, test.err, err)
		})
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
