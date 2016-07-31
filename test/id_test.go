package test

import (
	"testing"

	"github.com/flimzy/flashback-model"
)

func TestID(t *testing.T) {
	id := fbmodel.NewByteID("user", []byte("User Bob"))
	id2, err := fbmodel.ParseID(id.String())
	if err != nil {
		t.Fatalf("We can't even parse the IDs we generate: %s", err.Error())
	}
	if id.Type() != id2.Type() {
		t.Errorf("Types: %s != %s\n", id.Type(), id2.Type())
	}
	if id.Identity() != id2.Identity() {
		t.Errorf("ID: %x != %x\n", id.Identity(), id2.Identity())
	}
}

func TestID2(t *testing.T) {
	id := fbmodel.CreateID("user", []byte("User Bob"))
	id2, err := fbmodel.ParseID(id.String())
	if err != nil {
		t.Fatalf("We can't even parse the IDs we generate: %s", err.Error())
	}
	if id.Type() != id2.Type() {
		t.Errorf("Types: %s != %s\n", id.Type(), id2.Type())
	}
	if id.Identity() != id2.Identity() {
		t.Errorf("ID: %x != %x\n", id.Identity(), id2.Identity())
	}
}

type IDTest struct {
	Name  string
	ID    string
	Error string
}
