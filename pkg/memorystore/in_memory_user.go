package memorystore

import (
	"errors"
	"todolist/pkg/model"
)

type UserMap struct {
	user map[string]model.User
}

func NewUserMap() *UserMap {
	return &UserMap{
		user: make(map[string]model.User),
	}
}

func (u *UserMap) Create(user model.User) error {
	if _, ok := u.user[user.UID]; ok {
		return model.ErrObjectAlreadyExists
	}
	u.user[user.UID] = user
	return nil
}

func (u *UserMap) Login(username, password string) (bool, error) {
	for _, user := range u.user {
		if user.Username == username && user.Password == password {
			return true, nil
		}
	}
	return false, errors.New("user name and password are not valid")
}
