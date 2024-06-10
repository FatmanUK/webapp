package main

import (
	"log"
	"os"
//	"fatgo/mktls"
)

func Exists(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func createTls() {
	// path must exist
	key_is := (nil == Exists(conf.tls.key))
	crt_is := (nil == Exists(conf.tls.crt))
	if crt_is && key_is {
		return
	}
	if !key_is {
		//mktls.CreateKey(conf.tls.key)
	}
	//mktls.CreateCrt(conf.tls.key, conf.tls.crt)
}

var conf = &Config{file: configFile}

func main() {
	if nil != Exists(conf.file) {
		log.Output(0, "Creating default config, run again.")
		makeDefaults()
	} else {
		conf.load()
		createTls()
		createRoutes()
		err := run(conf.tls.key, conf.tls.crt)
		if err != nil {
			log.Fatal(err)
		}
	}
}
