package main

import (
	"os"
	"log"
)

type View struct {
	Config *JsonConfig
	Page *Page
	User User
}

func setDefaults(c *JsonConfig) {
	// TODO: automate this?
	c.SetString("web.appname", APPNAME)
	c.SetString("web.port", PORT)
	c.SetString("web.home", HOME)
	c.SetString("db.content", DATADIR + "/content.db")
	c.SetString("db.users", DATADIR + "/users.db")
	c.SetString("web.tls.key", CONFDIR + "/tls/tls.key")
	c.SetString("web.tls.crt", CONFDIR + "/tls/tls.crt")
	c.SetString("auth.keys_dir", DATADIR + "/keys")
	c.SetString("web.static_dir", DATADIR + "/static")
	c.SetString("web.template_dir", DATADIR + "/templates")
	c.SetString("web.timeouts.expiry_h", COOKIE_EXPIRY_HOURS)
	c.SetString("web.timeouts.idle_h", COOKIE_IDLE_TIMEOUT_HOURS)
	if ! fileExists(CONFIG_FILE) {
		mode := os.O_CREATE | os.O_WRONLY
		f, err := os.OpenFile(CONFIG_FILE, mode, os.ModePerm)
		defer f.Close()
		if err != nil {
			panic(errWriteConf)
		}
		c.Save(f)
	}
}

func main() {
	buildcmd := b64decode(BUILD_COMMAND_B64)
	log.Output(1, "Compile args: " + buildcmd)

	f, err := os.Open(CONFIG_FILE)
	defer f.Close()
	if err != nil {
		panic(errReadConf)
	}
	c := (&JsonConfig{}).Init(CONFIG_FILE)
	setDefaults(c)
	c.Load(f)

	(&Template{}).Init(c.GetString("web.template_dir"))
	(&User{}).OpenDatabase(c.GetString("db.users"))
	(&Page{}).OpenDatabase(c.GetString("db.content"))

	(&TlsHandler{}).GenerateSelfSignedCertificate(
		c.GetString("web.tls.key"),
		c.GetString("web.tls.crt"),
	)

	w := &WebRouter{ Config: c }
	err = w.Run()

	if err != nil {
		log.Fatal(err)
	}
}
