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

type handlerfn func(http.ResponseWriter, *http.Request, string)

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
		fn := makeHandler(re.debugHandler)
		http.HandleFunc("/debug/", fn)
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

func makeHandler(fn handlerfn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func (re *WebRouter) defaultHandler(
		w http.ResponseWriter,
		r *http.Request) {
	url := "/view/" + re.Config.GetString("web.home")
	http.Redirect(w, r, url, http.StatusFound)
}

func (re *WebRouter) viewHandler(
		w http.ResponseWriter,
		r *http.Request,
		title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session, re.Config)
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
		url := "/edit/" + title
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	v := &View{Config: re.Config, Page: p, User: *user}
	renderTemplate(w, "view", v)
}

func (re *WebRouter) editHandler(
		w http.ResponseWriter,
		r *http.Request,
		title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session, re.Config)
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

func (re *WebRouter) saveHandler(
		w http.ResponseWriter,
		r *http.Request,
		title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session, re.Config)
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

func (re *WebRouter) debugHandler(
		w http.ResponseWriter,
		r *http.Request,
		title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session, re.Config)
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

func (re *WebRouter) actionLogin(
		w http.ResponseWriter,
		user *User,
		page *Page) string {
	user.Login()
	re.SetCookie(w, 1, user.Session)
	log.Output(1, "Login attempt.")
	page.Title = "Login"
	return "userLogin"
}

func (re *WebRouter) actionLogin2(
		r *http.Request,
		user *User,
		page *Page) string {
	page.Title = "Access Denied"
	name := r.PostForm["User"][0]
	keydir := re.Config.GetString("auth.keys_dir")
	keyfile := keydir + "/" + name + ".asc"
	pubkey := loadTextFile(keyfile)
	template := "userLoginFailed"
	if isVerifiedPgpClearSignature(r.PostForm, user, pubkey) {
		template = "userWelcome"
		page.Title = "Hello"
		user.Authorise(name, re.Config)
		log.Output(1, "Login successful.")
	} else {
		log.Output(1, "Login failed.")
	}
	return template
}

func (re *WebRouter) actionLogout(
		w http.ResponseWriter,
		user *User,
		page *Page) string {
	user.Logout()
	re.SetCookie(w, -1, "")
	log.Output(1, "Logout by user")
	page.Title = "Logout"
	return "userLogout"
}

func (re *WebRouter) actionManage(page *Page) string {
	log.Output(1, "User management attempt")
	page.Title = "Manage Users"
	return "userManage"
}

func (re *WebRouter) actionCreate(
		r *http.Request,
		page *Page) string {
	name := r.PostForm["User"][0]
	keydir := re.Config.GetString("auth.keys_dir")
	keyfile := keydir + "/" + name + ".asc"
	saveTextFile(keyfile, r.PostForm["Datum"][0], 0600)
	(&User{
		Name: name,
		Nick: r.PostForm["Nick"][0],
	}).Create(r.PostForm["AddGroup"])
	log.Output(1, "User created")
	page.Title = "User Created"
	page.Body = []byte(name)
	return "userCreate"
}

func (re *WebRouter) actionDelete(
		r *http.Request,
		page *Page) string {
	name := r.PostForm["User"][0]
	keydir := re.Config.GetString("auth.keys_dir")
	deleteFile(keydir + "/" + name + ".asc")
	(&User{
		Name: name,
	}).Delete()
	page.Title = "User Deleted"
	page.Body = []byte(name)
	log.Output(1, "User delete")
	return "userDelete"
}

func (re *WebRouter) userHandler(
		w http.ResponseWriter,
		r *http.Request,
		a string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session, re.Config)
	template := "userDefault"
	p := Page{Title: "User Default"}
	if r.Method == "POST" {
		r.ParseForm()
	}
	switch a {
		case "login": {
			template = re.actionLogin(w, user, &p)
		}
		case "login2": {
			template = re.actionLogin2(r, user, &p)
		}
		case "logout": {
			template = re.actionLogout(w, user, &p)
		}
		case "manage": {
			if ! user.IsGroupMember("stewards") {
				denyUnauthorised(w, r)
				return
			}
			template = re.actionManage(&p)
		}
		case "create": {
			if ! user.IsGroupMember("stewards") {
				denyUnauthorised(w, r)
				return
			}
			template = re.actionCreate(r, &p)
		}
		case "delete": {
			if ! user.IsGroupMember("stewards") {
				denyUnauthorised(w, r)
				return
			}
			template = re.actionDelete(r, &p)
		}
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

func (re *WebRouter) SetCookie(
		w http.ResponseWriter,
		i int,
		session string) {
	timeout_h := re.Config.GetInt("web.timeouts.expiry_h")
	offset_s := time.Duration(i * timeout_h) * time.Hour
	cookie := &http.Cookie{
		Name: "session_token",
		Value: session,
		MaxAge: int(offset_s),
		Path: "/",
		Expires: time.Now().Add(offset_s),
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
