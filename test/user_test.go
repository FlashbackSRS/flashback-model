package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/flashback-model"
)

var frozenUser []byte = []byte(`{"_id":"user-nRHQJKEAQEWlt58cz5bMnw==","type":"user","password":"","salt":"","userType":"","username":"mrsmith"}`)

func TestUser(t *testing.T) {
	u, err := testUser()
	if err != nil {
		t.Fatalf("Error unfreezing test user: %s", err)
	}
	stringsEqual(t, "Username", u.Username, "mrsmith")
	stringsEqual(t, "ID", u.ID.String(), "user-nRHQJKEAQEWlt58cz5bMnw==")
	stringsEqual(t, "Type", u.Type(), "user")
	stringsEqual(t, "Identity", u.Identity(), "nRHQJKEAQEWlt58cz5bMnw==")
	output, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Error marshaling user: %s", err)
	}
	if equal, err := jsonDeepEqual(frozenUser, output); err != nil {
		t.Errorf("Unable to compare JSON output: %s\n", err)
	} else if !equal {
		t.Errorf("Unexpected JSON output\n     Got: %s\nExpected: %s\n", string(output), frozenUser)
	}
}

func TestNewUser(t *testing.T) {
	u, err := fbmodel.NewUser("9d11d024-a100-4045-a5b7-9f1ccf96cc9f")
	if err != nil {
		t.Errorf("Error creating user: %s\n", err)
	}
	stringsEqual(t, "ID", u.ID.String(), "user-nRHQJKEAQEWlt58cz5bMnw==")
	stringsEqual(t, "Type", u.Type(), "user")
	u.Username = "mrsmith"
	output, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Error marshaling new user: %s", err)
	}
	if equal, err := jsonDeepEqual(frozenUser, output); err != nil {
		t.Errorf("Unable to compare new user JSON output: %s\n", err)
	} else if !equal {
		t.Errorf("Unexpected JSON output\n     Got: %s\nExpected: %s\n", string(output), frozenUser)
	}
}

type UserTest struct {
	Name  string
	JSON  string
	Error string
}

func testUser() (*fbmodel.User, error) {
	u := &fbmodel.User{}
	err := json.Unmarshal([]byte(frozenUser), u)
	return u, err
}
