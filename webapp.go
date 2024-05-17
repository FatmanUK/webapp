package main

// ❯ botan keygen --algo=RSA --params=4096 >./tls.key
// ❯ botan gen_self_signed ./tls.key localhost >./tls.crt

import (
	"log"
)

// Global vars
// TODO: set up config class
var TLS_KEY = "tls/tls.key"
var TLS_CRT = "tls/tls.crt"
var webRoot = "/home/adam.richardson/webapp"

func main() {
	createRoutes()
	err := run()
	initStaticTemplates()
	if err != nil {
		log.Fatal(err)
	}
}
