package main

// ❯ botan keygen --algo=RSA --params=4096 >./tls.key
// ❯ botan gen_self_signed ./tls.key localhost >./tls.crt

import (
	"log"
)

func main() {
	createRoutes()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
