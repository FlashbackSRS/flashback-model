package fb

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/pborman/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		uuid     uuid.UUID
		username string
		expected *User
		err      string
	}{
		{
			name: "no UUID",
			err:  "invalid user id: id required",
		},
		{
			name: "no username",
			uuid: []byte("bob"),
			err:  "username required",
		},
		{
			name:     "valid",
			uuid:     []byte("bob"),
			username: "bob",
			expected: &User{
				ID:       DbID{docType: "user", id: []byte("bob")},
				uuid:     []byte("bob"),
				Username: "bob",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewUser(test.uuid, test.username)
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, result); err != nil {
				t.Error(d)
			}
		})
	}
}
