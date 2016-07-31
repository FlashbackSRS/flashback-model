package test

import (
	"encoding/json"
	"testing"

	"github.com/pborman/uuid"

	"github.com/flimzy/flashback-model"
	. "github.com/flimzy/flashback-model/test/util"
)

var frozenUser []byte = []byte(`{"_id":"user-nRHQJKEAQEWlt58cz5bMnw==","type":"user","password":"","salt":"","userType":"","username":"mrsmith"}`)

func TestUser(t *testing.T) {
	u, err := testUser()
	if err != nil {
		t.Fatalf("Error unfreezing test user: %s", err)
	}
	StringsEqual(t, "Username", u.Username, "mrsmith")
	StringsEqual(t, "ID", u.ID.String(), "user-nRHQJKEAQEWlt58cz5bMnw==")
	StringsEqual(t, "Type", u.Type(), "user")
	StringsEqual(t, "Identity", u.Identity(), "nRHQJKEAQEWlt58cz5bMnw==")
	output, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Error marshaling user: %s", err)
	}
	JSONDeepEqual(t, "Frozen user", frozenUser, output)
}

func TestNewUser(t *testing.T) {
	u, err := fbmodel.NewUser(uuid.Parse("9d11d024-a100-4045-a5b7-9f1ccf96cc9f"), "mrsmith")
	if err != nil {
		t.Errorf("Error creating user: %s\n", err)
	}
	StringsEqual(t, "ID", u.ID.String(), "user-nRHQJKEAQEWlt58cz5bMnw==")
	StringsEqual(t, "Type", u.Type(), "user")
	output, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Error marshaling new user: %s", err)
	}
	JSONDeepEqual(t, "New user", frozenUser, output)
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
