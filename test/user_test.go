package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"
	"github.com/pborman/uuid"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenUser = []byte(`
{
    "type": "user",
    "_id": "user-tui5ajfbabaeljnxt4om7fwmt4",
    "username": "mrsmith",
    "password": "",
    "salt": ""
}
`)

func TestNewUser(t *testing.T) {
	require := require.New(t)
	u, err := fb.NewUser(uuid.Parse("9d11d024-a100-4045-a5b7-9f1ccf96cc9f"), "mrsmith")
	require.Nil(err, "Error creating user: %s", err)
	require.Equal("user-tui5ajfbabaeljnxt4om7fwmt4", u.ID.String(), "User ID not as expected")
	require.Equal("user", u.ID.Type(), "User type not as expected")
	require.MarshalsToJSON(frozenUser, u, "New User")

	u2, err := testUser()
	require.Nil(err, "Error thawing user: %s", err)
	require.MarshalsToJSON(frozenUser, u2, "Thawed User")

	require.DeepEqual(u, u2, "Thawed vs. Created Users")
}

func testUser() (*fb.User, error) {
	u := &fb.User{}
	err := json.Unmarshal([]byte(frozenUser), u)
	return u, err
}
