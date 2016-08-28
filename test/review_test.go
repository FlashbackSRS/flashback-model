package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenReview = []byte(`
{
    "cardID": "VGVzdCBOb3Rl.0",
    "timestamp": null,
    "ease": 0,
    "interval": null,
    "previousInterval": null,
    "srsFactor": 0,
    "reviewTime": null,
    "reviewType": 0
}
`)

func TestReview(t *testing.T) {
	r, err := fb.NewReview("VGVzdCBOb3Rl.0")
	if err != nil {
		t.Fatalf("Error creating review: %s\n", err)
	}
	JSONDeepEqual(t, "Create Review", Marshal(t, "Create Review", r), frozenReview)

	r2 := &fb.Review{}
	if err := json.Unmarshal(frozenReview, r2); err != nil {
		t.Fatalf("Error thawing review: %s\n", err)
	}
	JSONDeepEqual(t, "Thawed Review", Marshal(t, "Tha Review", r2), frozenReview)

	if !reflect.DeepEqual(r, r2) {
		PrintDiff(r2, r)
		t.Fatal("Thawed and created Reviews don't match")
	}
}
