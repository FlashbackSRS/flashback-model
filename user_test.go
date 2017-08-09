package fb

import (
	"testing"

	"github.com/flimzy/diff"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected *User
		err      string
	}{
		{
			name: "no id",
			err:  "id required",
		},
		{
			name: "valid",
			id:   "user-mjxwe",
			expected: &User{
				ID: "user-mjxwe",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewUser(test.id)
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
		ID: "user-aaaaaaaaabaabaaaaaaaaaaaaa",
	}
	if d := diff.Interface(expected, u); d != "" {
		t.Error(d)
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
			name: "no id",
			user: &User{},
			err:  "id required",
		},
		{
			name: "null fields",
			user: &User{
				ID:       "user-mjxwe",
				Salt:     "salty",
				Password: "abc123",
			},
			expected: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123"
            }`,
		},
		{
			name: "all fields",
			user: &User{
				ID:       "user-mjxwe",
				Salt:     "salty",
				Password: "abc123",
				FullName: "Bob",
				Email:    "bob@bob.com",
			},
			expected: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123",
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

func TestUserUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *User
		err      string
	}{
		{
			name:  "invalid json",
			input: "invalid json",
			err:   "failed to unmarshal user: invalid character 'i' looking for beginning of value",
		},
		{
			name:  "wrong type",
			input: `{"type":"chicken"}`,
			err:   "Invalid document type for user",
		},
		{
			name:  "fails validation",
			input: `{"_id":"deck-mjxwe", "type":"user"}`,
			err:   "incorrect doc type",
		},
		{
			name: "null fields",
			input: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123"
            }`,
			expected: &User{
				ID:       "user-mjxwe",
				Salt:     "salty",
				Password: "abc123",
			},
		},
		{
			name: "all fields",
			input: `{
                "_id":      "user-mjxwe",
                "type":     "user",
                "salt":     "salty",
                "password": "abc123",
                "email":    "bob@bob.com",
                "fullname": "Bob"
            }`,
			expected: &User{
				ID:       "user-mjxwe",
				Salt:     "salty",
				Password: "abc123",
				Email:    "bob@bob.com",
				FullName: "Bob",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &User{}
			err := result.UnmarshalJSON([]byte(test.input))
			checkErr(t, test.err, err)
			if err != nil {
				return
			}
			if d := diff.Interface(test.expected, result); d != "" {
				t.Error(d)
			}
		})
	}
}

func TestUserValidate(t *testing.T) {
	tests := []validationTest{
		{
			name: "no ID",
			v:    &User{},
			err:  "id required",
		},
		{
			name: "invalid doctype",
			v:    &User{ID: "chicken-mjxwe"},
			err:  "incorrect doc type",
		},
		{
			name: "wrong doctype",
			v:    &User{ID: "bundle-mjxwe"},
			err:  "incorrect doc type",
		},
		{
			name: "valid",
			v:    &User{ID: "user-mjxwe"},
		},
	}
	testValidation(t, tests)
}
