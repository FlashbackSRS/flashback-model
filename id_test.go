package fb

import (
	"bytes"
	"testing"
)

func TestID(t *testing.T) {
	id, err := NewID("note", []byte("Test Note"))
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

func TestHexID(t *testing.T) {
	id, err := NewHexID("user", []byte("Test User"))
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
