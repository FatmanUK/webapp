package main

import (
	"os"
	"errors"
	"encoding/base64"
)

func b64decode(b string) string {
	d, e := base64.StdEncoding.DecodeString(b)
	if e != nil {
		panic(e.Error())
	}
	return string(d)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return ! errors.Is(err, os.ErrNotExist)
}

func loadTextFile(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}
