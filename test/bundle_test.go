package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/flimzy/flashback-model"
)

var frozenBundle []byte = []byte(`{"type":"bundle","_id":"bundle-dGVzdCBidW5kbGU=","owner":"nRHQJKEAQEWlt58cz5bMnw==","name":"Test Bundle","description":"A bundle for testing"}`)

func oink() {
}

func TestNewBundle(t *testing.T) {
	u, _ := testUser()
	b := fbmodel.NewBundle([]byte("test bundle"), u)
	name := "Test Bundle"
	b.Name = &name
	descr := "A bundle for testing"
	b.Description = &descr
	fmt.Printf("%v\n", b)
	stringsEqual(t, "Bundle ID", b.ID.String(), "bundle-dGVzdCBidW5kbGU=")
	jsonDeepEqual(t, "New Bundle", marshal(t, "New bundle", b), frozenBundle)
}

func TestThawBundle(t *testing.T) {
	b := &fbmodel.Bundle{}
	if err := json.Unmarshal(frozenBundle, b); err != nil {
		t.Fatalf("Error thawing bundle: %s", err)
	}
	jsonDeepEqual(t, "Thawed Bundle", marshal(t, "Thawed bundle", b), frozenBundle)
}
