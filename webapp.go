package main

import (
	"log"
	"os"
//	"fatgo/mktls"
)

// Global vars
// TODO: set up config class
var TLS_KEY = "tls/tls.key"
var TLS_CRT = "tls/tls.crt"
var webRoot = "/home/adam.richardson/webapp"

func Exists(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func createTls() {
	// path must exist
	key_is := (nil == Exists(TLS_KEY))
	crt_is := (nil == Exists(TLS_CRT))
	if crt_is && key_is {
		return
	}
	if !key_is {
		//mktls.CreateKey(TLS_KEY)
	}
	//mktls.CreateCrt(TLS_KEY, TLS_CRT)
}

func main() {
	createTls()
	createRoutes()
	initStaticTemplates()
	err := run(TLS_KEY, TLS_CRT)
	if err != nil {
		log.Fatal(err)
	}
}
