package user

import (
	"github.com/pborman/uuid"

	"github.com/flimzy/flashback-model"
)

type User struct {
	model.BaseDoc
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Salt     string  `json:"salt"`
	UserType string  `json:"userType"`
	FullName *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
	uuid     uuid.UUID
}

func Random() *User {
	return &User{
		uuid: uuid.NewRandom(),
	}
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}
