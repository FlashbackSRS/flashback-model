package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenB64ID []byte = []byte(`"note-0VGVzdCBOb3Rl"`)
var frozenHexID []byte = []byte(`"user-1546573742055736572"`)

func TestB64ID(t *testing.T) {
	id, err := fb.NewByteID("note", []byte("Test Note"), fb.Base64ID)
	if err != nil {
		t.Fatalf("Error creating B64 ID: %s\n", err)
	}
	JSONDeepEqual(t, "Create B64 ID", Marshal(t, "Create ID1", id), frozenB64ID)

	id2 := fb.ID{}
	if err := json.Unmarshal(frozenB64ID, &id2); err != nil {
		t.Fatalf("Error thawing B64 ID: %s", err)
	}
	JSONDeepEqual(t, "Thawed B64 ID", Marshal(t, "Thaw B64 ID", id2), frozenB64ID)

	if !reflect.DeepEqual(id, id2) {
		PrintDiff(id2, id)
		t.Fatalf("Thawed and created B64 IDs don't match")
	}
}

func TestHexID(t *testing.T) {
	id, err := fb.NewByteID("user", []byte("Test User"), fb.HexID)
	if err != nil {
		t.Fatalf("Error creating Hex ID: %s\n", err)
	}
	JSONDeepEqual(t, "Create Hex ID", Marshal(t, "Create ID1", id), frozenHexID)

	id2 := fb.ID{}
	if err := json.Unmarshal(frozenHexID, &id2); err != nil {
		t.Fatalf("Error thawing Hex ID: %s", err)
	}
	JSONDeepEqual(t, "Thawed Hex ID", Marshal(t, "Thaw Hex ID", id2), frozenHexID)

	if !reflect.DeepEqual(id, id2) {
		PrintDiff(id2, id)
		t.Fatalf("Thawed and created Hex IDs don't match")
	}
}

func TestID(t *testing.T) {
	id, err := fb.NewByteID("user", []byte("User Bob"), fb.HexID)
	if err != nil {
		t.Fatalf("Error creating user: %s\n", err)
	}
	id2, err := fb.ParseID(id.String())
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
	id, err := fb.CreateID("user", []byte("User Bob"))
	if err != nil {
		t.Fatalf("Error creating user: %s\n", err)
	}
	id2, err := fb.ParseID(id.String())
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
