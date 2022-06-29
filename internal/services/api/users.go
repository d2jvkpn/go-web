package api

import (
	"fmt"
	"sync"

	. "github.com/d2jvkpn/go-web/pkg/resp"
)

type UsersData struct {
	m     *sync.RWMutex
	users map[string]string
}

func NewUsersData() UsersData {
	ud := UsersData{
		m:     new(sync.RWMutex),
		users: make(map[string]string, 100),
	}

	// ud.users["admin"] = "admin"
	return ud
}

func (ud *UsersData) Exists(name string) (ok bool) {
	ud.m.RLock()
	defer ud.m.RUnlock()

	_, ok = ud.users[name]
	return ok
}

func (ud *UsersData) Register(user, password string) *HttpError {
	if !_UserRE.Match([]byte(user)) {
		msg := "invalid username"
		return ErrInvalidParameter(fmt.Errorf(msg), msg)
	}
	if !_PasswordRE.Match([]byte(password)) {
		msg := "invalid password"
		return ErrInvalidParameter(fmt.Errorf(msg), msg)
	}

	ud.m.Lock()
	defer ud.m.Unlock()

	if _, ok := ud.users[user]; ok {
		msg := "user is already exists"
		return ErrConflict(fmt.Errorf(msg+": "+user), msg)
	}

	ud.users[user] = password
	return nil
}

// Verify password before unregister
func (ud UsersData) Unregister(name string) {
	delete(ud.users, name)
}

func (ud UsersData) Verify(name, password string) *HttpError {
	if !_UserRE.Match([]byte(name)) || !_PasswordRE.Match([]byte(password)) {
		msg := "invalid username or password"
		return ErrInvalidParameter(fmt.Errorf(msg), msg)
	}

	ud.m.RLock()
	defer ud.m.RUnlock()

	if pw, ok := ud.users[name]; !ok {
		return ErrUnauthorized(fmt.Errorf("user doesn't exists"), "wrong username or password")
	} else if pw != password {
		return ErrUnauthorized(fmt.Errorf("password doesn't match"), "wrong username or password")
	}

	return nil
}
