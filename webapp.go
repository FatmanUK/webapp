package main

import (
	"log"
	"os"
//	"fatgo/mktls"
	"encoding/base64"
)

func fileExists(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func createTls() {
	// path must exist
	key_is := (nil == fileExists(c.GetString("tls.key")))
	crt_is := (nil == fileExists(c.GetString("tls.crt")))
	if crt_is && key_is {
		return
	}
	if !key_is {
		//mktls.CreateKey(conf.tls.key)
	}
	//mktls.CreateCrt(conf.tls.key, conf.tls.crt)
}

var c = &JsonConfig{
	CONFDIR + "/webapp.cfg",
	make(map[string]string)}

// must set these with build options
var APPNAME string
var HOME string
var PORT string
var CONFDIR string
var DATADIR string

var BUILD_MODE string
var BUILD_COMMAND_B64 string
var BUILD_COMMAND string

func defaults() {
	c.SetString("web.appname", APPNAME)
	c.SetString("web.port", PORT)
	c.SetString("web.home", HOME)
	c.SetString("db.content", DATADIR + "/content.db")
	c.SetString("db.users", DATADIR + "/users.db")
	c.SetString("tls.key", CONFDIR + "/tls/tls.key")
	c.SetString("tls.crt", CONFDIR + "/tls/tls.crt")
	c.SetString("keys_dir", DATADIR + "/keys")
	c.SetString("static_dir", DATADIR + "/static")
	if ! c.FileExists() {
		c.Save()
	}
}

func main() {
	defaults()
	c.Load()

	decoded, err := base64.StdEncoding.DecodeString(BUILD_COMMAND_B64)
	if err != nil {
		panic(err.Error())
	}
	BUILD_COMMAND = string(decoded)
	log.Output(0, "Compile args: " + BUILD_COMMAND)

	(&Template{}).Init()
	openDatabase()
	createTls()
	createRoutes()
	err = run()
	if err != nil {
		log.Fatal(err)
	}
}
