package main

type TlsHandler struct {
}

func (*TlsHandler) GenerateSelfSignedCertificate(
		key string,
		crt string) {
	// path must exist
	key_is := fileExists(key)
	crt_is := fileExists(crt)
	if crt_is && key_is {
		return
	}
	if !key_is {
		//mktls.CreateKey(key)
	}
	//mktls.CreateCrt(key, crt)
}
