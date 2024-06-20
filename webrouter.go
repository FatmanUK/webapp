package main

import (
	"net/http"
)

type WebRouter struct {
	Config *JsonConfig
}

type View struct {
	Page *Page
	User User
}

func (re *WebRouter) CreateRoutes() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/user/", makeHandler(userHandler))
	if BUILD_MODE == "Debug" {
		http.HandleFunc("/debug/", makeHandler(debugHandler))
	}
	dir := re.Config.GetString("web.static_dir")
	fs := http.FileServer(http.Dir(dir))
	t := http.StripPrefix("/static/", fs)
	http.Handle("/static/", t)
	http.HandleFunc("/", handler)
}

func (re *WebRouter) Run(c *JsonConfig) error {
	port := re.Config.GetString("web.port")
	key := re.Config.GetString("tls.key")
	cert := re.Config.GetString("tls.crt")
	return http.ListenAndServeTLS(":" + port, cert, key, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "/view/" + c.GetString("web.home")
	http.Redirect(w, r, url, http.StatusFound)
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	p, err := loadPage(title)
	if err != nil {
		if user.IsGroupMember("authors") {
			denyNotFound(w, r)
			return
		}
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", &View{Page: p, User: *user})
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	if user.IsGroupMember("authors") {
		denyUnauthorised(w, r)
		return
	}
	log.Output(1, "Editing " + title + ".")
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", &View{Page: p, User: *user})
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	if user.IsGroupMember("authors") {
		denyUnauthorised(w, r)
		return
	}
	log.Output(1, "Saving " + title + ".")
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func debugHandler(w http.ResponseWriter, r *http.Request, title string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	log.Output(1, "Debug tool accessed.")
	template := "debug"
	p := Page{Title: "Debug"}
	// TODO: improve this?
	output := User{}.Debug()
	output += c.debugOutput()
	output += Page{}.Debug()
	output += View{}.Debug()
	p.Body = []byte(output)
	renderTemplate(w, template, &View{Page: &p, User: *user})
}

func userHandler(w http.ResponseWriter, r *http.Request, a string) {
	var session string = ""
	cookie, _ := ReadSessionToken(r)
	if cookie != nil {
		session = cookie.Value
	}
	user := UserFromSessionToken(session)
	template := "userDefault"
	p := Page{Title: "User Default"}
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
		r.ParseForm()
		name := r.PostForm["User"][0]
		pubkey := loadTextFile(c.GetString("keys_dir") + "/" + name + ".asc")
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
	renderTemplate(w, template, &View{Page: &p, User: *user})
}

// for header.html
func (re *View) GetAppname() string {
	return c.GetString("web.appname")
}

// for header.html
func (re *View) GetIconname() string {
	return strings.ToLower(c.GetString("web.appname"))
}

func (re View) Debug() string {
	output := `
## View
___`
	return output
}

