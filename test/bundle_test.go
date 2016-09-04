package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/FlashbackSRS/flashback-model"
	. "github.com/FlashbackSRS/flashback-model/test/util"
)

var frozenBundle = []byte(`
{
    "type": "bundle",
    "_id": "bundle-krsxg5baij2w4zdmmu",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "owner": "tui5ajfbabaeljnxt4om7fwmt4",
    "name": "Test Bundle",
    "description": "A bundle for testing"
}
`)

func TestNewBundle(t *testing.T) {
	u, _ := testUser()
	b, err := fb.NewBundle([]byte("Test Bundle"), u)
	if err != nil {
		t.Fatalf("Error creating new bundle: %s", err)
	}
	name := "Test Bundle"
	b.Name = &name
	b.Created = now
	b.Modified = now
	imp := now.AddDate(0, 0, 2)
	b.Imported = &imp
	descr := "A bundle for testing"
	b.Description = &descr
	StringsEqual(t, "Bundle ID", b.ID.String(), "bundle-krsxg5baij2w4zdmmu")
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
