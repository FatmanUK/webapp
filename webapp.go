package main

import (
	"log"
	"os"
	"github.com/FatmanUK/fatgo/mktls"
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
	(&User{}).OpenDatabase()
	openDatabase()
	createTls()
	createRoutes()
	err = run()
	if err != nil {
		log.Fatal(err)
	}
}
