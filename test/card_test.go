package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenCard []byte = []byte(`
{
    "type": "card",
    "_id": "card-0mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z"
}
`)

func TestCard(t *testing.T) {
	c, err := fb.NewCard("0mViuXQThMLoh1G1Nlc4d_E8kR8o", 0)
	if err != nil {
		t.Fatalf("Error creating card: %s", err)
	}
	c.Created = now
	c.Modified = now
	imp := now.AddDate(0, 0, 2)
	c.Imported = &imp
	JSONDeepEqual(t, "Create Card", Marshal(t, "Create Card", c), frozenCard)

	c2 := &fb.Card{}
	if err := json.Unmarshal(frozenCard, c2); err != nil {
		t.Fatalf("Error thawing card: %s", err)
	}
	JSONDeepEqual(t, "Thawed Card", Marshal(t, "Thaw Card", c2), frozenCard)

	if !reflect.DeepEqual(c, c2) {
		PrintDiff(c2, c)
		t.Fatal("Thawed and created Cards don't match")
	}
}

/*
var frozenExistingCard []byte = []byte(`
{
    "type": "card",
    "_id": "card-0mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z"
}
`)

var frozenMergedCard []byte = []byte(`
{
    "type": "card",
    "_id": "card-0mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z"
}
`)

func TestCardMergeImport(t *testing.T) {
	c := &fb.Card{}
	if err := json.Unmarshal(frozenCard, c); err != nil {
		t.Fatalf("Error thawing Card: %s", err)
	}
	e := &fb.Card{}
	if err := json.Unmarshal(frozenExistingCard, e); err != nil {
		t.Fatalf("Error thawing ExistingCard: %s", err)
	}
	changed, err := c.MergeImport(e)
	if err != nil {
		t.Fatalf("Error merging Card: %s\n", err)
	}
	if !changed {
		t.Fatalf("No change!")
	}
	JSONDeepEqual(t, "Merged Card", Marshal(t, "Merge Card", c), frozenMergedCard)
}
*/
