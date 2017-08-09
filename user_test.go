package fb

import (
	"bytes"
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

func TestNilUser(t *testing.T) {
	u := NilUser()
	expected := &User{
		ID:       DbID{docType: "user", id: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		uuid:     []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Username: "niluser",
	}
	if d := diff.Interface(expected, u); d != "" {
		t.Error(d)
	}
}

func TestUserUUID(t *testing.T) {
	expected := []byte("foo")
	u := &User{uuid: expected}
	if id := u.UUID(); !bytes.Equal(id, expected) {
		t.Errorf("Unexpected result: %v", id)
	}
}

func TestUserMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected string
		err      string
	}{
		{
			name: "null fields",
			user: &User{
				ID:       DbID{docType: "user", id: []byte("bob")},
				Username: "bob",
				Salt:     "salty",
				Password: "abc123",
			},
			expected: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123",
                "username": "bob"
            }`,
		},
		{
			name: "all fields",
			user: &User{
				ID:       DbID{docType: "user", id: []byte("bob")},
				Username: "bob",
				Salt:     "salty",
				Password: "abc123",
				FullName: func() *string { x := "Bob"; return &x }(),
				Email:    func() *string { x := "bob@bob.com"; return &x }(),
			},
			expected: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123",
                "username": "bob",
                "email":    "bob@bob.com",
                "fullname": "Bob"
            }`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.user.MarshalJSON()
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.JSON([]byte(test.expected), result); d != "" {
				t.Error(d)
			}
		})
	}
}
