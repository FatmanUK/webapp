package main

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
