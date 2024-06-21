package main

import (
	"os"
	"github.com/FatmanUK/fatgo/mktls"
)

type TlsHandler struct {
}

func (re *TlsHandler) GenerateSelfSignedCertificate(
		key string,
		crt string) {
	// path must exist
	key_is := fileExists(key)
	crt_is := fileExists(crt)

	if crt_is && key_is {
		return
	}

	k := (&mktls.TlsKey{}).GenerateKey()
	c := k.GenerateCertificate()

	mode := os.O_CREATE|os.O_WRONLY|os.O_TRUNC
	kf, err := os.OpenFile(key, mode, os.ModePerm)
	defer kf.Close()
	if err != nil {
		panic(errWriteKey)
	}
	kf.Write(k.PemBytes.Bytes())

	cf, err := os.OpenFile(crt, mode, os.ModePerm)
	defer cf.Close()
	if err != nil {
		panic(errWriteCert)
	}
	cf.Write(c.PemBytes.Bytes())
}
