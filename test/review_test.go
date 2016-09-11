package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenReview = []byte(`
{
    "cardID": "VGVzdCBOb3Rl.0",
    "timestamp": null
}
`)

func TestReview(t *testing.T) {
	require := require.New(t)
	r, err := fb.NewReview("VGVzdCBOb3Rl.0")
	require.Nil(err, "Error creating review: %s", err)

	require.MarshalsToJSON(frozenReview, r, "Create Review")

	r2 := &fb.Review{}
	err = json.Unmarshal(frozenReview, r2)
	require.Nil(err, "Error thawing review: %s", err)
	require.MarshalsToJSON(frozenReview, r2, "Thawed Review")

	require.DeepEqual(r, r2, "Thawed vs. Created Reviews")
}
