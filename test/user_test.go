package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/flimzy/flashback-model"
)

var frozenUser []byte = []byte(`{"_id":"user-9d11d024-a100-4045-a5b7-9f1ccf96cc9f","type":"user","password":"","salt":"","userType":"","name":"mrsmith"}`)

func TestUser(t *testing.T) {
	o := testUser()
	fmt.Printf("user = %v\n", o)
	output, err := json.Marshal(o)
	if err != nil {
		t.Errorf("Error marshaling user: %s", err)
	}
	if equal, err := jsonDeepEqual(frozenUser, output); err != nil {
		t.Errorf("Unable to compare JSON output: %s\n", err)
	} else if !equal {
		t.Errorf("Unexpected JSON output\n     Got: %s\nExpected: %s\n", string(output), frozenUser)
	}
	if id := o.ID(); id != "9d11d024-a100-4045-a5b7-9f1ccf96cc9f" {
		t.Errorf("Unexpected user id:\n     Got: %s\nExpected: %s\n", id, "9d11d024-a100-4045-a5b7-9f1ccf96cc9f")
	}
	if tp := o.Type(); tp != "user" {
		t.Errorf("Unexpected user type:\n     Got: %s\nExpected: %s\n", tp, "user")
	}
	fmt.Printf("Output = %s\n", string(output))
}

func TestNewUser(t *testing.T) {
	u, err := fbmodel.NewUser("9d11d024-a100-4045-a5b7-9f1ccf96cc9f")
	if err != nil {
		t.Errorf("Error creating user: %s\n", err)
	}
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

var UserTests []UserTest = []UserTest{
	UserTest{
		Name:  "Invalid UUID",
		JSON:  `{"_id":"user-9d11d024-a100-404x-a5b7-9f1ccf96cc9f","type":"user","password":"","salt":"","userType":"","name":"bobsmith"}`,
		Error: "Invalid user ID: 9d11d024-a100-404x-a5b7-9f1ccf96cc9f",
	},
	UserTest{
		Name:  "Invalid Type",
		JSON:  `{"_id":"oink-9d11d024-a100-4045-a5b7-9f1ccf96cc9f","type":"user","password":"","salt":"","userType":"","name":"Bob Smith"}`,
		Error: "Invalid document type: oink",
	},
}

func TestBrokenUsers(t *testing.T) {
	for _, test := range UserTests {
		u := &fbmodel.User{}
		err := json.Unmarshal([]byte(test.JSON), u)
		if err.Error() != test.Error {
			t.Errorf("[%s] Unexpected error.\n     Got: %s\nExpected: %s\n", test.Name, err.Error(), test.Error)
		}
	}
}

func testUser() *fbmodel.User {
	u := &fbmodel.User{}
	json.Unmarshal([]byte(frozenUser), u)
	return u
}
