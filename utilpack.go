package main

import (
	"os"
	"io/fs"
	"errors"
	"encoding/base64"
	"time"
)

func stringFromZuluTime(t *time.Time) string {
	u := []byte("- nil -")
	if t != nil {
		u, _ = t.UTC().MarshalText()
	}
	return string(u)
}

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

func saveTextFile(file string, content string, mode fs.FileMode) {
	err := os.WriteFile(file, []byte(content), mode)
	if err != nil {
		panic(err.Error())
	}
}

func loadTextFile(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}
