package main

import (
	"os"
	"github.com/FatmanUK/fatgo/mktls"
	"encoding/base64"
//	"log"
//	"errors"
)

var BUILD_COMMAND_B64 string
var BUILD_MODE string

var CONFIG_FILE string
var WEB_PORT string
var WEB_HOME string
var DB_PAGES string
var DB_USERS string
var TLS_KEY string
var TLS_CERT string
var AUTH_KEYS_DIR string

var errNoWriteConf = "Couldn't write config"

type TlsHandler struct {
}

func (*TlsHandler) GenerateSelfSignedCertificate(
		key string,
		crt string) {
	// path must exist
	key_path := c.GetString("tls.key")
	cert_path := c.GetString("tls.crt")
	key_is := (nil == fileExists(key_path))
	crt_is := (nil == fileExists(cert_path))
	if crt_is && key_is {
		return
	}
//	if !key_is {
//		mktls.CreateKey(conf.tls.key)
//	}
//	mktls.CreateCrt(conf.tls.key, conf.tls.crt)
	k := (&mktls.TlsKey{}).GenerateKey()
	c := k.GenerateCertificate()

	if !key_is || !crt_is {
		mode := os.O_CREATE|os.O_WRONLY|os.O_TRUNC
		kf, err := os.OpenFile(key_path, mode, os.ModePerm)
		defer kf.Close()
		if err != nil {
			panic("Couldn't write key file")
		}
		kf.Write(k.PemBytes.Bytes())

		cf, err := os.OpenFile(cert_path, mode, os.ModePerm)
		defer cf.Close()
		if err != nil {
			panic("Couldn't write cert file")
		}
		cf.Write(c.PemBytes.Bytes())
	}
}

//var c = &JsonConfig{
//		CONFDIR + "/webapp.cfg",
//		make(map[string]string),
//	}

// must set these with build options
var APPNAME string
var HOME string
var PORT string
var CONFDIR string
var DATADIR string

var BUILD_MODE string
var BUILD_COMMAND_B64 string
var BUILD_COMMAND string

type WebRouter struct {
}

func (re *WebRouter) CreateRoutes() {
}

func (re *WebRouter) Run() error {
	return nil
}

func setDefaults(c *JsonConfig) {
	c.SetString("web.appname", APPNAME)
	c.SetString("web.port", PORT)
	c.SetString("web.home", HOME)
	c.SetString("db.content", DATADIR + "/content.db")
	c.SetString("db.users", DATADIR + "/users.db")
	c.SetString("web.tls.key", CONFDIR + "/tls/tls.key")
	c.SetString("web.tls.crt", CONFDIR + "/tls/tls.crt")
	c.SetString("auth.keys_dir", DATADIR + "/keys")
	c.SetString("web.static_dir", DATADIR + "/static")

	if ! fileExists(c.File) {
		mode := os.O_CREATE | os.O_WRONLY
		f, err := os.OpenFile(c.File, mode, os.ModePerm)
		defer f.Close()
		if err != nil {
			panic(errNoWriteConf)
		}
		c.Save(f)
	}
}

func main() {
	f, err := os.Open(CONFIG_FILE)
	defer f.Close()
	if err != nil {
		panic(errNoReadConf)
	}
	c := (&JsonConfig{}).Init(CONFIG_FILE)
	setDefaults(c)
	c.Load(f)

	(&Template{}).Init()
	(&User{}).OpenDatabase()
	openDatabase()

	buildcmd := b64decode(BUILD_COMMAND_B64)
	log.Output(1, "Compile args: " + buildcmd)

	p := (&Page{}).OpenDatabase()
	u := (&User{}).OpenDatabase()

	(&TlsHandler{}).GenerateSelfSignedCertificate(
		c.GetString("web.tls.key"),
		c.GetString("web.tls.crt"),
	)

	w := &WebRouter{ Config: c }
	w.CreateRoutes()
	err = w.Run()

	if err != nil {
		log.Fatal(err)
	}
}
