package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenBundle []byte = []byte(`
{
    "type": "bundle",
    "_id": "bundle-VjMOV9J35iuH1lXdM_lgQPOYx9I=",
    "owner": "nRHQJKEAQEWlt58cz5bMnw==",
    "name": "Test Bundle",
    "description": "A bundle for testing"
}
`)

func TestNewBundle(t *testing.T) {
	u, _ := testUser()
	b, err := fb.NewBundle("VjMOV9J35iuH1lXdM_lgQPOYx9I=", u)
	if err != nil {
		t.Fatalf("Error creating new bundle: %s", err)
	}
	name := "Test Bundle"
	b.Name = &name
	descr := "A bundle for testing"
	b.Description = &descr
	StringsEqual(t, "Bundle ID", b.ID.String(), "bundle-VjMOV9J35iuH1lXdM_lgQPOYx9I=")
	JSONDeepEqual(t, "New Bundle", Marshal(t, "New bundle", b), frozenBundle)

	b2 := &fb.Bundle{}
	if err := json.Unmarshal(frozenBundle, b2); err != nil {
		t.Fatalf("Error thawing bundle: %s", err)
	}
	JSONDeepEqual(t, "Thawed Bundle", Marshal(t, "Thawed bundle", b2), frozenBundle)

	// We have to set the username explicitly for the next test to pass, as a simple unmarshaling
	// of a bundle doesn't know user details (nor should it)
	b2.Owner.Username = "mrsmith"
	if !reflect.DeepEqual(b, b2) {
		PrintDiff(b2, b)
		t.Fatalf("Thawed and created Bundles don't match")
	}
}
