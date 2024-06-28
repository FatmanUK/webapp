package main

import (
	"log"
	"net/http"
	"strings"
	"regexp"
	"time"
)

// See also:
// https://go.dev/doc/articles/wiki/

var validPath = regexp.MustCompile(pathRe)
var pathRe = "^/(edit|save|view|user|debug)/([a-zA-Z0-9]+)$"

type WebRouter struct {
	Config *JsonConfig
}

func (re *WebRouter) Run() error {
	// create routes
	http.HandleFunc("/view/", makeHandler(re.viewHandler))
	http.HandleFunc("/edit/", makeHandler(re.editHandler))
	http.HandleFunc("/save/", makeHandler(re.saveHandler))
	http.HandleFunc("/user/", makeHandler(re.userHandler))
	if BUILD_MODE == "Debug" {
		http.HandleFunc("/debug/", makeHandler(re.debugHandler))
	}
	dir := re.Config.GetString("web.static_dir")
	fs := http.FileServer(http.Dir(dir))
	t := http.StripPrefix("/static/", fs)
	http.Handle("/static/", t)
	http.HandleFunc("/", re.defaultHandler)

	// run server
	port := re.Config.GetString("web.port")
	key := re.Config.GetString("web.tls.key")
	cert := re.Config.GetString("web.tls.crt")
	return http.ListenAndServeTLS(":" + port, cert, key, nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func (re *WebRouter) defaultHandler(w http.ResponseWriter, r *http.Request) {
	url := "/view/" + re.Config.GetString("web.home")
	http.Redirect(w, r, url, http.StatusFound)
}

func (re *WebRouter) viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	t := time.Now().UTC()
	user.LastRequest = &t
	user.Save()
	p := &Page{}
	err := p.LoadPage(title)
	if err != nil {
		if ! user.IsGroupMember("authors") {
			denyNotFound(w, r)
			return
		}
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	v := &View{Config: re.Config, Page: p, User: *user}
	renderTemplate(w, "view", v)
}

func (re *WebRouter) editHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	if user.Name == "" {
		denyAuthReqd(w, r)
		return
	}
	if ! user.IsGroupMember("authors") {
		denyUnauthorised(w, r)
		return
	}
	log.Output(1, "Editing " + title + ".")

	p := &Page{}
	err := p.LoadPage(title)
	if err != nil {
		p.Title = title
	}
	v := &View{Config: re.Config, Page: p, User: *user}
	renderTemplate(w, "edit", v)
}

func (re *WebRouter) saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	if user.Name == "" {
		denyAuthReqd(w, r)
		return
	}
	if ! user.IsGroupMember("authors") {
		denyUnauthorised(w, r)
		return
	}
	log.Output(1, "Saving " + title + ".")
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.Save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func (re *WebRouter) debugHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	log.Output(1, "Debug tool accessed")
	template := "debug"
	p := Page{Title: "Debug"}
	// TODO: improve this?
	output := User{}.Debug()
	output += re.Config.Debug()
	output += (&p).Debug()
	output += View{}.Debug()
	p.Body = []byte(output)
	v := &View{Config: re.Config, Page: &p, User: *user}
	renderTemplate(w, template, v)
}

func (re *WebRouter) userHandler(w http.ResponseWriter, r *http.Request, a string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	keydir := re.Config.GetString("auth.keys_dir")
	user := UserFromSessionToken(session)
	template := "userDefault"
	p := Page{Title: "User Default"}
	if r.Method == "POST" {
		r.ParseForm()
	}
	if a == "login" {
		template = "userLogin"
		user.Login()
		SetCookie(w, 1, user.Session)
		p.Title = "Login"
		log.Output(1, "Login attempt.")
	}
	if a == "login2" {
		template = "userLoginFailed"
		p.Title = "Access Denied"
		name := r.PostForm["User"][0]
		keyfile := keydir + "/" + name + ".asc"
		pubkey := loadTextFile(keyfile)
		if isVerifiedPgpClearSignature(r.PostForm, user, pubkey) {
			template = "userWelcome"
			p.Title = "Hello"
			user.Authorise(name)
			log.Output(1, "Login successful.")
		} else {
			log.Output(1, "Login failed.")
		}
	}
	if a == "logout" {
		template = "userLogout"
		p.Title = "Logout"
		user.Logout()
		SetCookie(w, -1, "")
		log.Output(1, "Logout by user")
	}
	if a == "manage" {
		if ! user.IsGroupMember("stewards") {
			denyUnauthorised(w, r)
			return
		}
		template = "userManage"
		p.Title = "Manage Users"
		log.Output(1, "User management attempt")
	}
	if a == "create" {
		if ! user.IsGroupMember("stewards") {
			denyUnauthorised(w, r)
			return
		}
		template = "userCreate"
		p.Title = "User Created"
		name := r.PostForm["User"][0]
		keyfile := keydir + "/" + name + ".asc"
		saveTextFile(keyfile, r.PostForm["Datum"][0], 0600)
		(&User{
			Name: name,
			Nick: r.PostForm["Nick"][0],
		}).Create(r.PostForm["AddGroup"])
		p.Body = []byte(name)
		log.Output(1, "User created")
	}
	if a == "delete" {
		if ! user.IsGroupMember("stewards") {
			denyUnauthorised(w, r)
			return
		}
		template = "userDelete"
		p.Title = "User Deleted"
		name := r.PostForm["User"][0]
		deleteFile(keydir + "/" + name + ".asc")
		(&User{
			Name: name,
		}).Delete()
		p.Body = []byte(name)
		log.Output(1, "User delete")
	}
	v := &View{Config: re.Config, Page: &p, User: *user}
	renderTemplate(w, template, v)
}

// for header.html
func (re *View) GetAppname() string {
	return re.Config.GetString("web.appname")
}

// for header.html
func (re *View) GetIconname() string {
	return strings.ToLower(re.Config.GetString("web.appname"))
}

func (re View) Debug() string {
	output := `
## View
___`
	return output
}

func SetCookie(w http.ResponseWriter, i int, session string) {
	now := time.Now()
	offset := 24 * time.Duration(i) * time.Hour
	expiry := now.Add(offset)
	cookie := &http.Cookie{
		Name: "session_token",
		Value: session,
		MaxAge: i * 86400,
		Path: "/",
		Expires: expiry,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func ReadSessionToken(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("session_token")
}

func denyNotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusNotFound)
	log.Output(1, "Not found.")
}

func denyAuthReqd(w http.ResponseWriter, r *http.Request) {
	// The "unauthorized" status actually means "unauthenticated".
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
	http.Redirect(w, r, "/", http.StatusUnauthorized)
	log.Output(1, "Login required.")
}

func denyUnauthorised(w http.ResponseWriter, r *http.Request) {
	// We use "Forbidden" to mean "unauthorised".
	http.Redirect(w, r, "/", http.StatusForbidden)
	log.Output(1, "Not allowed.")
}
