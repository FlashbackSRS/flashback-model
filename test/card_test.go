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
    "_id": "card-0mViuXQThMLoh1G1Nlc4d_E8kR8o.0"
}
`)

func TestCard(t *testing.T) {
	c, err := fb.NewCard("0mViuXQThMLoh1G1Nlc4d_E8kR8o", 0)
	if err != nil {
		t.Fatalf("Error creating card: %s", err)
	}
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
