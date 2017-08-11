package fb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var nilUser = uuid.UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

// User repressents a user of Flashback
type User struct {
	ID        string    `json:"_id"`
	Rev       string    `json:"_rev,omitempty"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	FullName  string    `json:"fullname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
	LastLogin time.Time `json:"lastLogin,omitempty"`
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
	if u.Created.IsZero() {
		return errors.New("created time required")
	}
	if u.Modified.IsZero() {
		return errors.New("modified time required")
	}
	return nil
}

// NewUser returns a new User object, based on the provided UUID and username.
func NewUser(id string) (*User, error) {
	u := &User{
		ID:       id,
		Created:  now(),
		Modified: now(),
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

// NilUser returns a special user, whose UUID bits are all set to zero, to be
// used as a placeholder when the actual user isn't known.
func NilUser() *User {
	u, _ := NewUser(EncodeDBID("user", nilUser))
	return u
}

type userAlias User

// MarshalJSON implements the json.Marshaler interface for the User type.
func (u *User) MarshalJSON() ([]byte, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	doc := struct {
		userAlias
		LastLogin *time.Time `json:"lastLogin,omitempty"`
	}{
		userAlias: userAlias(*u),
	}
	if !u.LastLogin.IsZero() {
		doc.LastLogin = &u.LastLogin
	}
	return json.Marshal(doc)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the User type.
func (u *User) UnmarshalJSON(data []byte) error {
	doc := &userAlias{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal user")
	}
	*u = User(*doc)
	return u.Validate()
}
