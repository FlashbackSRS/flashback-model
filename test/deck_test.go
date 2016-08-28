package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "name": "Test Deck",
    "description": "Deck for testing",
    "cards": []
}
`)

func TestDecks(t *testing.T) {
	d, err := fb.NewDeck([]byte("Test Deck"))
	if err != nil {
		t.Fatalf("Error creating deck: %s", err)
	}
	name := "Test Deck"
	d.Name = &name
	descr := "Deck for testing"
	d.Description = &descr
	d.Created = now
	d.Modified = now
	imp := now.AddDate(0, 0, 2)
	d.Imported = &imp
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

var frozenExistingDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z",
    "name": "Test deck",
    "description": "Deck for testing",
    "cards": []
}
`)

var frozenMergedDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "name": "Test Deck",
    "description": "Deck for testing",
    "cards": []
}
`)

func TestDeckMergeImport(t *testing.T) {
	d := &fb.Deck{}
	if err := json.Unmarshal(frozenDeck, d); err != nil {
		t.Fatalf("Error thawing Deck: %s", err)
	}
	e := &fb.Deck{}
	if err := json.Unmarshal(frozenExistingDeck, e); err != nil {
		t.Fatalf("Error thawing ExistingDeck: %s", err)
	}
	changed, err := d.MergeImport(e)
	if err != nil {
		t.Fatalf("Error merging Deck: %s\n", err)
	}
	if !changed {
		t.Fatalf("No change in deck merge")
	}
	JSONDeepEqual(t, "Merged Deck", Marshal(t, "Merge Deck", d), frozenMergedDeck)
}
