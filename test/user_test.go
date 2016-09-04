package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pborman/uuid"

	"github.com/FlashbackSRS/flashback-model"
	. "github.com/FlashbackSRS/flashback-model/test/util"
)

var frozenUser = []byte(`
{
    "type": "user",
    "_id": "user-tui5ajfbabaeljnxt4om7fwmt4",
    "username": "mrsmith",
    "password": "",
    "salt": "",
    "userType": ""
}
`)

func TestNewUser(t *testing.T) {
	u, err := fb.NewUser(uuid.Parse("9d11d024-a100-4045-a5b7-9f1ccf96cc9f"), "mrsmith")
	if err != nil {
		t.Errorf("Error creating user: %s\n", err)
	}
	StringsEqual(t, "ID", u.ID.String(), "user-tui5ajfbabaeljnxt4om7fwmt4")
	StringsEqual(t, "Type", u.ID.Type(), "user")
	JSONDeepEqual(t, "New user", Marshal(t, "New User", u), frozenUser)

	u2, err := testUser()
	if err != nil {
		t.Fatalf("Error thawing user: %s", err)
	}
	JSONDeepEqual(t, "New user", Marshal(t, "New User", u2), frozenUser)

	if !reflect.DeepEqual(u, u2) {
		PrintDiff(u2, u)
		t.Fatalf("Thawed and created Users don't match")
	}
}

func testUser() (*fb.User, error) {
	u := &fb.User{}
	err := json.Unmarshal([]byte(frozenUser), u)
	return u, err
}
