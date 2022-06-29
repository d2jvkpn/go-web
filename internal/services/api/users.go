package api

import (
	"fmt"
	"sync"

	. "github.com/d2jvkpn/go-web/pkg/resp"

	"golang.org/x/crypto/bcrypt"
)

type UsersData struct {
	m     *sync.RWMutex
	users map[string][]byte
}

func NewUsersData() UsersData {
	ud := UsersData{
		m:     new(sync.RWMutex),
		users: make(map[string][]byte, 100),
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

	var (
		bts []byte
		err error
	)

	if _, ok := ud.users[user]; ok {
		msg := "user is already exists"
		return ErrConflict(fmt.Errorf(msg+": "+user), msg)
	}

	if bts, err = bcrypt.GenerateFromPassword([]byte(password), _HASH_Cost); err != nil {
		return ErrServerError(err)
	}

	ud.users[user] = bts
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

	var (
		ok  bool
		bts []byte
		err error
	)

	msg := "wrong username or password"
	if bts, ok = ud.users[name]; !ok {
		return ErrUnauthorized(fmt.Errorf("user doesn't exists"), msg)
	}
	if err = bcrypt.CompareHashAndPassword(bts, []byte(password)); err != nil {
		return ErrUnauthorized(err, msg)
	}

	return nil
}
