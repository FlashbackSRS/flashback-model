package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenDocID = []byte(`"note-VGVzdCBOb3Rl"`)
var frozenDbID = []byte(`"user-krsxg5bakvzwk4q"`)

func TestDbID(t *testing.T) {
	require := require.New(t)
	id, err := fb.NewDbID("user", []byte("Test User"))
	require.Nil(err, "Error creating DbID: %s", err)
	require.Equal("user-krsxg5bakvzwk4q", id.String(), "Stringified DbID not as expected")
	require.Equal("krsxg5bakvzwk4q", id.Identity(), "DbID identity not as expected")
	require.MarshalsToJSON(frozenDbID, id, "Create DbID")

	id2 := fb.DbID{}
	err = json.Unmarshal(frozenDbID, &id2)
	require.Nil(err, "Error thawing DbID: %s", err)
	require.MarshalsToJSON(frozenDbID, id2, "Thawed DbID")

	require.DeepEqual(id, id2, "Thawed vs. Created DbIDs")
}

func TestID(t *testing.T) {
	require := require.New(t)
	id, err := fb.NewDbID("user", []byte("User Bob"))
	require.Nil(err, "Error creating user: %s", err)
	id2, err := fb.ParseDbID(id.String())
	require.Nil(err, "We can't even parse the IDs we generate: %s", err)
	require.Equal(id.Type(), id2.Type(), "ID Type equality")
	require.Equal(id.Identity(), id2.Identity(), "ID Identity equality")
}

func TestID2(t *testing.T) {
	id, err := fb.NewDbID("user", []byte("User Bob"))
	if err != nil {
		t.Fatalf("Error creating user: %s\n", err)
	}
	id2, err := fb.ParseDbID(id.String())
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
