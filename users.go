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

var sessions map[string]*User = map[string]*User{}

func (re User) debugOutput() string {
	output := `
## Sessions`
	for k, _ := range sessions {
		output += `
    Session:      ` + k + `  
    User.Name:    ` + sessions[k].Name + `  
    User.Session: ` + sessions[k].Session + `  
    User.Nonce:   ` + sessions[k].Nonce
		for _, v := range sessions[k].Groups {
			output += `
    User.Group:   ` + v
		}
	}
	output += `
___`
	return output
}

func (re *User) authorise() {
	// TODO: do something with user
	// get groups
	re.Groups = []string{"authors"}
}

func userFromCookie(r *http.Request) (*User, error) {
	c, err := r.Cookie("session_token")
	if err == nil {
		return sessions[c.Value], nil
	}
	return nil, errors.New("There Is No Cookie")
}

func userLogin(session string, w http.ResponseWriter) *User {
	delete(sessions, session)
	session = betterguid.New()
	user := &User{Session: session, Nonce: betterguid.New()}
	for user.Session == user.Nonce {
		user.Nonce = betterguid.New()
	}
	sessions[session] = user
	expiry := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name: "session_token",
		Value:   session,
		Expires: expiry,
		Path: "/",
		SameSite: http.SameSiteLaxMode }
	http.SetCookie(w, cookie)
	return user
}
