package fb

import (
	"encoding/json"
	"strings"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var nilUser = uuid.UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

// User repressents a user of Flashback
type User struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	FullName string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
}

// Validate validates that all of the data in the user appears valid and self
// consistent. A nil return value means no errors were detected.
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("id required")
	}
	if !strings.HasPrefix(u.ID, "user-") {
		return errors.New("incorrect doc type")
	}
	return nil
}

// NewUser returns a new User object, based on the provided UUID and username.
func NewUser(id, username string) (*User, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	if username == "" {
		return nil, errors.New("username required")
	}
	return &User{
		ID:       id,
		Username: username,
	}, nil
}

// NilUser returns a special user, whose UUID bits are all set to zero, to be
// used as a placeholder when the actual user isn't known.
func NilUser() *User {
	u, _ := NewUser("user-aaaaaaaaabaabaaaaaaaaaaaaa", "niluser")
	return u
}

type userAlias User

type jsonUser struct {
	userAlias
	Type string `json:"type"`
}

// MarshalJSON implements the json.Marshaler interface for the User type.
func (u *User) MarshalJSON() ([]byte, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	doc := struct {
		jsonUser
	}{
		jsonUser: jsonUser{
			Type:      "user",
			userAlias: userAlias(*u),
		},
	}
	return json.Marshal(doc)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the User type.
func (u *User) UnmarshalJSON(data []byte) error {
	doc := jsonUser{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal user")
	}
	if doc.Type != "user" {
		return errors.New("Invalid document type for user")
	}
	*u = User(doc.userAlias)
	return u.Validate()
}
