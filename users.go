package main

import (
	"net/http"
	"errors"
	"time"
	"github.com/kjk/betterguid"
)

type User struct {
	Name string
	Groups []string
	Nonce string // for login, nil after one use
	Session string // session identifier
}

var users map[string]User // should have called this sessions

func userFromCookie(r *http.Request) (User, error) {
	c, err := r.Cookie("session_token")
	if err == nil {
		return users[c.Value], nil
	}
	return User{}, errors.New("There Is No Cookie")
}

func userLogin(user *User, w http.ResponseWriter) {
	delete(users, user.Session)

	session := betterguid.New()
	user = &User{Session: session, Nonce: betterguid.New()}

	users[session] = *user

	// write session cookie
	expiry := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name: "session_token",
		Value:   session,
		Expires: expiry }

	http.SetCookie(w, cookie)
}
