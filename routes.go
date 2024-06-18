package main

import (
	"net/http"
	"regexp"
	"slices"
)

type View struct {
	Page *Page
	User User
}

func (re View) debugOutput() string {
	output := `
## View
___`
	return output
}

var pathRe = "^/(edit|save|view|user|debug)/([a-zA-Z0-9]+)$"
var validPath = regexp.MustCompile(pathRe)

func createRoutes() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/user/", makeHandler(userHandler))
	if BUILD_MODE == "Debug" {
		http.HandleFunc("/debug/", makeHandler(debugHandler))
	}
	fs := http.FileServer(http.Dir("static"))
	t := http.StripPrefix("/static/", fs)
	http.Handle("/static/", t)
	http.HandleFunc("/", handler)
}

// https://go.dev/doc/articles/wiki/

func run() error {
	port := c.GetString("web.port")
	key := c.GetString("tls.key")
	cert := c.GetString("tls.crt")
	return http.ListenAndServeTLS(":" + port, cert, key, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "/view/" + c.GetString("web.first_page")
	http.Redirect(w, r, url, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	user, _ := userFromCookie(r)
	if user == nil {
		user = &User{}
	}
	p, err := loadPage(title)
	if err != nil {
		if ! slices.Contains(user.Groups, "authors") {
			denyNotFound(w, r)
			return
		}
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", &View{Page: p, User: *user})
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	user, _ := userFromCookie(r)
	if user == nil {
		user = &User{}
	}
	if user.Name == "" {
		denyAuthReqd(w, r)
		return
	}
	if ! slices.Contains(user.Groups, "authors") {
		denyUnauthorised(w, r)
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", &View{Page: p, User: *user})
}

func denyNotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusNotFound)
}

func denyAuthReqd(w http.ResponseWriter, r *http.Request) {
	// The "unauthorized" status actually means "unauthenticated".
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
	http.Redirect(w, r, "/", http.StatusUnauthorized)
}

func denyUnauthorised(w http.ResponseWriter, r *http.Request) {
	// We use "Forbidden" to mean "unauthorised".
	http.Redirect(w, r, "/", http.StatusForbidden)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	user, _ := userFromCookie(r)
	if user == nil {
		user = &User{}
	}
	if user.Name == "" {
		denyAuthReqd(w, r)
		return
	}
	if ! slices.Contains(user.Groups, "authors") {
		denyUnauthorised(w, r)
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func debugHandler(w http.ResponseWriter, r *http.Request, title string) {
	user, _ := userFromCookie(r)
	if user == nil {
		user = &User{}
	}
	template := "debug"
	p := Page{Title: "Debug"}
	// TODO: improve this?
	output := User{}.debugOutput()
	output += c.debugOutput()
	output += Page{}.debugOutput()
	output += View{}.debugOutput()
	p.Body = []byte(output)
	renderTemplate(w, template, &View{Page: &p, User: *user})
}

func userHandler(w http.ResponseWriter, r *http.Request, a string) {
	user, _ := userFromCookie(r)
	if user == nil {
		user = &User{}
	}
	template := "userDefault"
	p := Page{Title: "User Default"}
	if a == "login" {
		template = "userLogin"
		user = userLogin(user.Session, w)
		p.Title = "Login"
	}
	if a == "login2" {
		template = "userLoginFailed"
		p.Title = "Access Denied"
		r.ParseForm()
		pubkey := loadTextFile(c.GetString("keys_dir") + "/" + r.PostForm["User"][0] + ".asc")
		if isVerifiedPgpClearSignature(r.PostForm, user, pubkey) {
			template = "userWelcome"
			p.Title = "Hello"
		}
	}
	if a == "logout" {
		template = "userLogout"
		p.Title = "Logout"
		userLogout(user.Session, w)
		log.Output(1, "Logout by user")
	}
	renderTemplate(w, template, &View{Page: &p, User: *user})
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
