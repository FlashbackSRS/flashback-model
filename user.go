package fbmodel

import (
	"encoding/json"
	"errors"

	"github.com/pborman/uuid"
)

type userTypeID struct {
	typeID
}

func isValidID(id string) bool {
	return uuid.Parse(id) != nil
}

type User struct {
	userTypeID
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
	ID       string  `json:"_id"`
	Rev      *string `json:"_rev,omitempty"`
	Username string  `json:"name"`
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
	u.userTypeID.set("user", id)
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
		ID:       u.DocID(),
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
	if err := u.userTypeID.parse(doc.ID); err != nil {
		return err
	}
	if !isValidID(u.docID) {
		return errors.New("Invalid user ID: " + u.docID)
	}
	u.Rev = doc.Rev
	u.Username = doc.Username
	u.Salt = doc.Salt
	u.FullName = doc.FullName
	u.Email = doc.Email

	return nil
}
