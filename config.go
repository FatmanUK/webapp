package main

import (
	"os"
	"dreamtrack.net/fatgo/serialisers"
)

type ConfigWeb struct {
	root string
	first_page string
}

type ConfigTls struct {
	key string
	crt string
}

type Config struct {
	file string
	web ConfigWeb
	tls ConfigTls
}

// can change these with build options
var configFile = "webapp.cfg"
var defaultWebRoot = "data"
var defaultTlsKey = "tls/tls.key"
var defaultTlsCrt = "tls/tls.crt"
var defaultFirstPage = "FrontPage"

// hmm, set these in build options?
func makeDefaults() {
	conf.web.root = defaultWebRoot
	conf.web.first_page = defaultFirstPage
	conf.tls.key = defaultTlsKey
	conf.tls.crt = defaultTlsCrt
	conf.save()
}

func (re *Config) serialise(s serialisers.Serialiser) {
	s.IoS(&re.web.root)
	s.IoS(&re.web.first_page)
	s.IoS(&re.tls.key)
	s.IoS(&re.tls.crt)
}

// strange bug...
/*
$ go run .  # makeDefaults triggered : size, save, load
Serialised:  data
Serialised:  tls/tls.key
Serialised:  tls/tls.crt
Serialised:  data
Serialised:  tls/tls.key
Serialised:  tls/tls.crt
Serialised:  datadata
Serialised:  tls/tls.keytls/tls.key
Serialised:  tls/tls.crttls/tls.crt
2024/05/17 16:06:07 open tls/tls.crttls/tls.crt: no such file or directory
exit status 1
$ go run .  # makeDefaults not triggered : just load
Serialised:  data
Serialised:  tls/tls.key
Serialised:  tls/tls.crt

# Seems like something is getting reused in Loader from Saver?
# Need to force a garbage collection when they're destroyed?
*/

func (re *Config) load() error {
	data, err := os.ReadFile(conf.file)
	if err != nil {
		return err
	}
	re.serialise(&serialisers.Loader{Array: &data})
	return nil
}

func (re *Config) save() error {
	var buf_size uint64 = 0
	re.serialise(&serialisers.Sizer{&buf_size})
	buf := make([]byte, buf_size)
	re.serialise(&serialisers.Saver{Array: &buf})
	return os.WriteFile(conf.file, buf, 0600)
}
