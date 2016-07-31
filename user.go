package fbmodel

import (
	"encoding/json"
	"errors"
	"fmt"
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

func NewUser(id string) (*User, error) {
	u := &User{}
	userUUID := uuid.Parse(id)
	if userUUID == nil {
		return nil, errors.New("Invalid user ID: " + id)
	}
	u.uuid = userUUID
	u.ID = NewID("user", userUUID)
	return u, nil
}

func RandomUser() *User {
	u, _ := NewUser(uuid.New())
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
		fmt.Printf("json error: %s\n", err)
		return err
	}
	fmt.Printf("doc = %s\n", doc)
	u.ID = doc.ID
	fmt.Printf("doc.ID = %#v\n", doc.ID)
	u.Rev = doc.Rev
	u.Username = doc.Username
	u.Salt = doc.Salt
	u.FullName = doc.FullName
	u.Email = doc.Email

	return nil
}
