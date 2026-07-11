package acl

import (
	"crypto/sha256"
	"encoding/hex"
)

type Storage struct {
	data map[string]*User
}

func NewStorage() *Storage {
	return &Storage{
		data: map[string]*User{
			"default" : {Name: "default", NoPass: true},
		},
	}
}

type User struct {
	Name string
	Password string
	NoPass bool
}

func (u *User) SetPassword(pass string) {
	hash := sha256.Sum256([]byte(pass))
	u.Password = hex.EncodeToString(hash[:])
	u.NoPass = false
}