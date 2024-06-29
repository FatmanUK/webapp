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
var timers map[string]*time.Timer = map[string]*time.Timer{}

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

func UserFromSessionToken(session string, c *JsonConfig) *User {
	u, exists := sessions[session]
	if ! exists {
		return &User{}
	}
	u.ResetIdleTimeout(c)
	return u
}

func (re *User) Authorise(name string, c *JsonConfig) {
	dur_h := time.Duration(c.GetInt("web.timeouts.idle_h"))
	timers[re.Session] = re.CreateLogoutTimer(time.Hour * dur_h)

	dur_h = time.Duration(c.GetInt("web.timeouts.expiry_h"))
	re.CreateLogoutTimer(time.Hour * dur_h)

	result := userDb.Where("name = ?", name).First(re)
	err := result.Error

	t := time.Now().UTC()
	re.LastLogin = &t
	re.Save()
	if err != nil {
		panic(errLoadUser)
	}
}

func (re User) Debug() string {
	// time.Timer can't report time left. Missing feature.
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
	output += `___`
	return output
}

func (re *User) CreateLogoutTimer(t time.Duration) *time.Timer {
	return time.AfterFunc(t, re.Logout)
}

func (re *User) ResetIdleTimeout(c* JsonConfig) {
	dur_h := time.Duration(c.GetInt("web.timeouts.idle_h"))
	t, exists := timers[re.Session]
	if exists {
		t.Reset(time.Hour * dur_h)
	}
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
