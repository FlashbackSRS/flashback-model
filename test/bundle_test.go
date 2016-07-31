package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenBundle []byte = []byte(`{"type":"bundle","_id":"bundle-VjMOV9J35iuH1lXdM_lgQPOYx9I=","owner":"nRHQJKEAQEWlt58cz5bMnw==","name":"Test Bundle","description":"A bundle for testing"}`)

func oink() {
}

func TestNewBundle(t *testing.T) {
	u, _ := testUser()
	b, err := fbmodel.NewBundle("VjMOV9J35iuH1lXdM_lgQPOYx9I=", u)
	if err != nil {
		t.Fatalf("Error creating new bundle: %s", err)
	}
	name := "Test Bundle"
	b.Name = &name
	descr := "A bundle for testing"
	b.Description = &descr
	fmt.Printf("%v\n", b)
	StringsEqual(t, "Bundle ID", b.ID.String(), "bundle-VjMOV9J35iuH1lXdM_lgQPOYx9I=")
	JSONDeepEqual(t, "New Bundle", Marshal(t, "New bundle", b), frozenBundle)
}

func TestThawBundle(t *testing.T) {
	b := &fbmodel.Bundle{}
	if err := json.Unmarshal(frozenBundle, b); err != nil {
		t.Fatalf("Error thawing bundle: %s", err)
	}
	JSONDeepEqual(t, "Thawed Bundle", Marshal(t, "Thawed bundle", b), frozenBundle)
}
