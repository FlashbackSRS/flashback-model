package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenDeck []byte = []byte(`
{
    "type": "deck",
    "_id": "deck-AO1yee9FPLVtU3h0M5pcYy3AOTQ=",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "name": "Test Deck",
    "description": "Deck for testing",
	"cards": []
}
`)

func TestDecks(t *testing.T) {
	d, err := fb.NewDeck("AO1yee9FPLVtU3h0M5pcYy3AOTQ=")
	if err != nil {
		t.Fatalf("Error creating deck: %s", err)
	}
	name := "Test Deck"
	d.Name = &name
	descr := "Deck for testing"
	d.Description = &descr
	d.Created = &now
	d.Modified = &now
	JSONDeepEqual(t, "Create Deck", Marshal(t, "Create Deck", d), frozenDeck)

	d2 := &fb.Deck{}
	if err := json.Unmarshal(frozenDeck, d2); err != nil {
		t.Fatalf("Error thawing deck: %s", err)
	}
	JSONDeepEqual(t, "Thawed Deck", Marshal(t, "Thaw Deck", d2), frozenDeck)

	if !reflect.DeepEqual(d, d2) {
		PrintDiff(d2, d)
		t.Fatal("Thawed and created Decks don't match")
	}
}
