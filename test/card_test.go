package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenCard = []byte(`
{
    "type": "card",
    "_id": "card-krsxg5baij2w4zdmmu.mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "model": "theme-VGVzdCBUaGVtZQ/0",
    "due": "2017-01-01",
    "interval": 50
}
`)

func TestCard(t *testing.T) {
	require := require.New(t)
	u, _ := testUser()
	b, _ := fb.NewBundle([]byte("Test Bundle"), u)
	c, err := fb.NewCard("theme-VGVzdCBUaGVtZQ", 0, "card-"+b.ID.Identity()+".mViuXQThMLoh1G1Nlc4d_E8kR8o.0")
	require.Nil(err, "Error creating card: %s", err)

	c.Created = now
	c.Modified = now
	c.Imported = now.AddDate(0, 0, 2)
	due, _ := fb.ParseDue("2017-01-01")
	c.Due = due
	ivl, _ := fb.ParseInterval("50d")
	c.Interval = ivl
	require.MarshalsToJSON(frozenCard, c, "Created Card")

	c2 := &fb.Card{}
	err = json.Unmarshal(frozenCard, c2)
	require.Nil(err, "Error thawing card: %s", err)
	require.MarshalsToJSON(frozenCard, c2, "Thawed Card")

	require.DeepEqual(c, c2, "Thawed vs Created Cards")
}

/*
var frozenExistingCard = []byte(`
{
    "type": "card",
    "_id": "card-mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z"
}
`)

var frozenMergedCard = []byte(`
{
    "type": "card",
    "_id": "card-mViuXQThMLoh1G1Nlc4d_E8kR8o.0",
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
