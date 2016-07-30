package fbmodel

import (
	"encoding/json"

// "fmt"
	"github.com/pborman/uuid"
)

type User struct {
	doc
	userDoc
	uuid     uuid.UUID
}

type userDoc struct {
	baseDoc
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Salt     string  `json:"salt"`
	UserType string  `json:"userType"`
	FullName *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func RandomUser() *User {
	u := &User{}
	userUUID := uuid.NewRandom()
	u.doc = NewDoc("user", userUUID.String())
	u.uuid = userUUID
	return u
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}

type jsonUser struct {
	*jsonDoc
	*User
}

// func (j *jsonUser) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data,j)
// // 	panic("foo")
// // 	return nil
// }

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonUser{
		jsonDoc: u.jsonDoc(),
		User: u,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	doc := &userDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	u.userDoc = *doc
	return nil
}
