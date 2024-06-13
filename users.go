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

// should have called this sessions
var users map[string]User = map[string]User{}

func (re User) debugOutput() string {
	output := `
## Sessions
  
`
	for k, _ := range sessions {
		output += `
    Session:      ` + k + `  
    User.Name:    ` + sessions[k].Name + `  
    User.Session: ` + sessions[k].Session + `  
    User.Nonce:   ` + sessions[k].Nonce /*+ `  
    User.Groups:  ` + sessions[k].Groups + `  
`*/
	}
	output += `
___`
	return output
}

func userFromCookie(r *http.Request) (*User, error) {
	c, err := r.Cookie("session_token")
	if err == nil {
		return users[c.Value], nil
	}
	return User{}, errors.New("There Is No Cookie")
}

func userLogin(session string, w http.ResponseWriter) User {
	delete(users, session)

	session = betterguid.New()
	user := User{Session: session, Nonce: betterguid.New()}

	users[session] = user

	// write session cookie
	expiry := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name: "session_token",
		Value:   session,
		Expires: expiry }

	http.SetCookie(w, cookie)
	return user
}