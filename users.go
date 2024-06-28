package main

import (
	"fmt"
	"strings"
	"time"
	"gorm.io/gorm"
	"github.com/glebarez/sqlite" // pure Go?
	"github.com/kjk/betterguid"
)

type User struct {
	Name string
	Nick string
	Groups string
	Created *time.Time
	LastLogin *time.Time
	LastRequest *time.Time
	Nonce string `gorm:"-"`
	Session string `gorm:"-"`
}

var sessions map[string]*User = map[string]*User{}

var userDb *gorm.DB

var debugFormat string = `
    Session:           %s  
    User.Name:         %s  
    User.Nick:         %s  
    User.Groups:       %s  
    User.Created:      %s  
    User.LastLogin:    %s  
    User.LastRequest:  %s  
    User.Session:      %s  
    User.Nonce:        %s  `

func (*User) OpenDatabase(dbfile string) {
	d, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	userDb = d
	userDb.AutoMigrate(&User{})
}

func UserFromSessionToken(session string) *User {
	u, exists := sessions[session]
	if ! exists {
		return &User{}
	}
	return u
}

func (re *User) Authorise(name string) {
	err := re.Load(name)
	t := time.Now().UTC()
	re.LastLogin = &t
	re.Save()
	if err != nil {
		panic(errLoadUser)
	}
}

func (re User) Debug() string {
	output := `
## Sessions`
	for k, _ := range sessions {
		created := sessions[k].Created
		last_login := sessions[k].LastLogin
		last_request := sessions[k].LastRequest
		output += fmt.Sprintf(
			debugFormat,
			k,
			sessions[k].Name,
			sessions[k].Nick,
			sessions[k].Groups,
			stringFromZuluTime(created),
			stringFromZuluTime(last_login),
			stringFromZuluTime(last_request),
			sessions[k].Session,
			sessions[k].Nonce,
		)
	}
	output += `
___`
	return output
}

func (re *User) Load(name string) error {
	result := userDb.Where("name = ?", name).First(re)
	return result.Error
}

func (re *User) Save() {
	if re.Name != "" {
		model := userDb.Model(re)
		model.Where("name = ?", re.Name).Updates(*re)
	}
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

func (re *User) Create(groups []string) {
	t := time.Now().UTC()
	re.Created = &t
	re.Groups = strings.Join(groups, ";")
	userDb.Create(re)
}

func (re *User) Delete() {
	if re.Name != "" {
		model := userDb.Model(re)
		model.Where("name = ?", re.Name).Delete(*re)
	}
}
