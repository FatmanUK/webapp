package main

import (
	"log"
	"os"
//	"fatgo/mktls"
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

// must set these with build options
var configFile string
var defaultWebRoot string
var defaultTlsKey string
var defaultTlsCrt string
var defaultFirstPage string

var c = &JsonConfig{
	configFile,
	make(map[string]string)}

func defaults() {
	c.SetString("web.root", defaultWebRoot)
	c.SetString("web.first_page", defaultFirstPage)
	c.SetString("tls.key", defaultTlsKey)
	c.SetString("tls.crt", defaultTlsCrt)
	if ! c.FileExists() {
		c.Save()
	}
}

func main() {
	defaults()
	c.Load()

	createTls()
	createRoutes()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
