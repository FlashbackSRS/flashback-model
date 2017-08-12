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
		ID:       "user-aaaaaaaaabaabaaaaaaaaaaaaa",
		Created:  now(),
		Modified: now(),
	}
	if d := diff.Interface(expected, u); d != nil {
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
				Created:  now(),
				Modified: now(),
			},
			expected: `{
			    "type":     "user",
				"_id":      "user-mjxwe",
				"salt":     "salty",
				"password": "abc123",
				"created":  "2017-01-01T00:00:00Z",
				"modified": "2017-01-01T00:00:00Z"
            }`,
		},
		{
			name: "all fields",
			user: &User{
				ID:        "user-mjxwe",
				Salt:      "salty",
				Password:  "abc123",
				FullName:  "Bob",
				Email:     "bob@bob.com",
				Created:   now(),
				Modified:  now(),
				LastLogin: now(),
			},
			expected: `{
				"type":     "user",
				"_id":       "user-mjxwe",
				"salt":      "salty",
				"password":  "abc123",
				"email":     "bob@bob.com",
				"fullname":  "Bob",
				"created":   "2017-01-01T00:00:00Z",
				"modified":  "2017-01-01T00:00:00Z",
				"lastLogin": "2017-01-01T00:00:00Z"
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
			if d := diff.JSON([]byte(test.expected), result); d != nil {
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
			name:  "fails validation",
			input: `{"_id":"deck-mjxwe"}`,
			err:   "incorrect doc type",
		},
		{
			name: "null fields",
			input: `{
				"_id":      "user-mjxwe",
				"salt":     "salty",
				"password": "abc123",
				"created":  "2017-01-01T00:00:00Z",
				"modified": "2017-01-01T00:00:00Z"
            }`,
			expected: &User{
				ID:       "user-mjxwe",
				Salt:     "salty",
				Password: "abc123",
				Created:  now(),
				Modified: now(),
			},
		},
		{
			name: "all fields",
			input: `{
				"_id":       "user-mjxwe",
				"salt":      "salty",
				"password":  "abc123",
				"email":     "bob@bob.com",
				"fullname":  "Bob",
				"created":   "2017-01-01T00:00:00Z",
				"modified":  "2017-01-01T00:00:00Z",
				"lastLogin": "2017-01-01T00:00:00Z"
            }`,
			expected: &User{
				ID:        "user-mjxwe",
				Salt:      "salty",
				Password:  "abc123",
				Email:     "bob@bob.com",
				FullName:  "Bob",
				Created:   now(),
				Modified:  now(),
				LastLogin: now(),
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
			if d := diff.Interface(test.expected, result); d != nil {
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
			name: "no created time",
			v:    &User{ID: "user-mzxw6"},
			err:  "created time required",
		},
		{
			name: "no modified time",
			v:    &User{ID: "user-mzxw6", Created: now()},
			err:  "modified time required",
		},
		{
			name: "valid",
			v:    &User{ID: "user-mjxwe", Created: now(), Modified: now()},
		},
	}
	testValidation(t, tests)
}
