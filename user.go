package fbmodel

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/pborman/uuid"
)

func isValidID(id string) bool {
	return uuid.Parse(id) != nil
}

type User struct {
	ID
	uuid     uuid.UUID
	Rev      *string
	Username string
	Password string
	Salt     string
	UserType string
	FullName *string
	Email    *string
}

type userDoc struct {
	Type     string  `json:"type"`
	ID       ID      `json:"_id"`
	Rev      *string `json:"_rev,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Salt     string  `json:"salt"`
	UserType string  `json:"userType"`
	FullName *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
}

// func CreateUser(username string) (*User, error) {
// 	u := &User{}
// 	u.ID = NewID("user", uuid.NewRandom())
// 	return u, nil
// }

func NewUser(id uuid.UUID, username string) (*User, error) {
	u := &User{}
	u.uuid = id
	u.ID = NewByteID("user", id)
	u.Username = username
	return u, nil
}

func NewUserStub(id string) (*User, error) {
	data, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}
	userUUID := uuid.UUID(data)
	return NewUser(userUUID, "")
}

func RandomUser() *User {
	u, _ := NewUser(uuid.NewRandom(), "randomuser")
	return u
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(userDoc{
		Type:     "user",
		ID:       u.ID,
		Rev:      u.Rev,
		Username: u.Username,
		Password: u.Password,
		Salt:     u.Salt,
		UserType: u.UserType,
		FullName: u.FullName,
		Email:    u.Email,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	doc := userDoc{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	if doc.Type != "user" {
		return errors.New("Invalid document type for user")
	}
	u.ID = doc.ID
	u.Rev = doc.Rev
	u.Username = doc.Username
	u.Salt = doc.Salt
	u.FullName = doc.FullName
	u.Email = doc.Email

	return nil
}

func (u *User) Fleshened() bool {
	return u.Username != ""
}
