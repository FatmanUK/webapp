package main

import (
	"io/fs"
	"regexp"
	"strings"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// gpg --export --armor f@d.n >keys/f@d.n.asc

type GpgPublicKey struct {
	armoured_data []byte
}

func (re *GpgPublicKey) LoadFileS(filename string) {
	re.armoured_data = loadTextFile(filename)
}

func (re *GpgPublicKey) LoadBlobS(blob string) {
	re.armoured_data = []byte(blob)
}

func (re *GpgPublicKey) SaveSI(file string, mode fs.FileMode) {
	saveTextFile(file, re.armoured_data, mode)
}

func (re *GpgPublicKey) bNonceMatchSS(
		nonce string,
		datum string) bool {
	// I'm going to rule that the file path *needs* to be "-"
	// here. It's beyond the scope of this project to determine
	// what is and isn't a valid file path, which is no simple
	// thing. Easier to insist on the stdin "pseudo-file", which
	// is a single hyphen.
	var nonceRe = "- ([^\n]*)"
	r := regexp.MustCompile(nonceRe)
	m := r.FindStringSubmatch(datum)
	if m == nil {
		return false
	}
	return (strings.TrimSpace(m[1]) == nonce)
}

func (re *GpgPublicKey) bIsGoodHashS(hash string) bool {
	switch hash {
		case "SHA256", "SHA512": {
			return true
		}
	}
	return false
}

func (re *GpgPublicKey) bIsValidClearSignatureAs(
		form map[string][]string) bool {
	var bsmRe = "-----BEGIN PGP SIGNED MESSAGE-----"
	var bpsRe = "-----BEGIN PGP SIGNATURE-----"
	var epsRe = "-----END PGP SIGNATURE-----"
	var hashRe = "Hash: ([A-z0-9-_]*)"
	ok := true
	msg := form["Datum"][0]
	r := regexp.MustCompile(hashRe)
	m := r.FindStringSubmatch(msg)
	ok = ok && (m != nil)
	if ok {
		ok = ok && re.bIsGoodHashS(m[1])
	}
	regexes := []string{bsmRe, bpsRe, epsRe}
	for _, v := range regexes {
		r = regexp.MustCompile(v)
		m = r.FindStringSubmatch(msg)
		ok = ok && (m != nil)
	}
	return ok
}

func (re *GpgPublicKey) bIsGoodClearSignatureMssU(
		r map[string][]string,
		user *User) bool {
	if ! re.bIsValidClearSignatureAs(r) {
		return false
	}
	name := r["User"][0]
	if name == "" {
		return false
	}
	nonce := user.Nonce
	user.Nonce = "" // only use once
	datum := r["Datum"][0]
	if ! re.bNonceMatchSS(nonce, datum) {
		return false
	}
	// decode plaintext from cleartext signature
	pt, err := helper.VerifyCleartextMessageArmored(
			string(re.armoured_data),
			datum,
			crypto.GetUnixTime())
	if pt != nonce || err != nil {
		return false
	}
	return true
}
