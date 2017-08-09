package fb

import (
	"encoding/json"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var nilUser = uuid.UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

// User repressents a user of Flashback
type User struct {
	ID       DbID
	uuid     uuid.UUID
	Rev      *string
	Username string
	Password string
	Salt     string
	FullName *string
	Email    *string
}

type userDoc struct {
	Type     string  `json:"type"`
	ID       DbID    `json:"_id"`
	Rev      *string `json:"_rev,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Salt     string  `json:"salt"`
	FullName *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
}

// Validate validates that all of the data in the user appears valid and self
// consistent. A nil return value means no errors were detected.
func (u *userDoc) Validate() error {
	if len(u.ID.id) == 0 {
		return errors.New("id required")
	}
	if u.ID.docType != "user" {
		return errors.New("incorrect doc type")
	}
	return nil
}

// NewUser returns a new User object, based on the provided UUID and username.
func NewUser(id uuid.UUID, username string) (*User, error) {
	uid, err := NewDbID("user", id)
	if err != nil {
		return nil, errors.Wrap(err, "invalid user id")
	}
	if username == "" {
		return nil, errors.New("username required")
	}
	return &User{
		ID:       uid,
		uuid:     id,
		Username: username,
	}, nil
}

// NilUser returns a special user, whose UUID bits are all set to zero, to be
// used as a placeholder when the actual user isn't known.
func NilUser() *User {
	u, _ := NewUser(nilUser, "niluser")
	return u
}

// UUID returns the UUID underlying the User object.
func (u *User) UUID() uuid.UUID {
	return u.uuid
}

// MarshalJSON implements the json.Marshaler interface for the User type.
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(userDoc{
		Type:     "user",
		ID:       u.ID,
		Rev:      u.Rev,
		Username: u.Username,
		Password: u.Password,
		Salt:     u.Salt,
		FullName: u.FullName,
		Email:    u.Email,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for the User type.
func (u *User) UnmarshalJSON(data []byte) error {
	doc := userDoc{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal user")
	}
	if doc.Type != "user" {
		return errors.New("Invalid document type for user")
	}
	u.ID = doc.ID
	u.uuid = u.ID.RawID()
	u.Rev = doc.Rev
	u.Username = doc.Username
	u.Password = doc.Password
	u.Salt = doc.Salt
	u.FullName = doc.FullName
	u.Email = doc.Email

	return nil
}
