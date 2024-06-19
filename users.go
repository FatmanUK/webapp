package main

import (
	"gorm.io/gorm"
	"github.com/glebarez/sqlite" // pure Go?
	"strings"
	"github.com/kjk/betterguid"
)

var errNoUserFound = "user data not found"

type User struct {
	Name string
	Nick string
	Groups string
	Nonce string `gorm:"-"`
	Session string `gorm:"-"`
}

var sessions map[string]*User = map[string]*User{}

var userDb *gorm.DB

func UserFromSessionToken(session string) (*User) {
	u, exists := sessions[session]
	if ! exists {
		return &User{}
	}
	return u
}

func (*User) OpenDatabase() {
	dbfile := "users.db"
	d, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	userDb = d
	userDb.AutoMigrate(&User{})
}

func (re *User) Authorise(name string) {
	err := re.Load(name)
	// anything else?
	if err != nil {
		panic(errNoUserFound)
	}
}

func (re User) Debug() string {
	output := `
## Sessions`
	for k, _ := range sessions {
		output += `
    Session:      ` + k + `  
    User.Name:    ` + sessions[k].Name + `  
    User.Nick:    ` + sessions[k].Nick + `  
    User.Groups:  ` + sessions[k].Groups + `  
    User.Session: ` + sessions[k].Session + `  
    User.Nonce:   ` + sessions[k].Nonce
	}
	output += `
___`
	return output
}

func (re *User) Load(name string) error {
	result := userDb.Where("name = ?", name).First(re)
	return result.Error
}

func (re *User) IsGroupMember(group string) bool {
	v := strings.Split(re.Groups, ";")
	for _, g := range v {
		if g == group {
			return true
		}
	}
	return false
}

func (re *User) Logout() {
	delete(sessions, re.Session)
}

func (re *User) Login() {
	delete(sessions, re.Session)
	re.Session = betterguid.New()
	re.Nonce = betterguid.New()
	// bug in betterguid, sometimes produces the same value
	for re.Session == re.Nonce {
		re.Nonce = betterguid.New()
	}
	sessions[re.Session] = re

}
