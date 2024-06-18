package main

import (
	"os"
	"net/url"
	"regexp"
	"strings"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// gpg --export --armor f@d.n >keys/f@d.n.asc

// I'm going to rule that the file path *needs* to be "-" here. It's
// beyond the scope of this project to determine what is and isn't a
// valid file path, which is no simple thing. Easier to insist on the
// stdin "pseudo-file", which is a single hyphen.
var nonceRe = "- ([^\n]*)"

func makeMatch(msg string, re string) ([]string, bool) {
	r := regexp.MustCompile(re)
	m := r.FindStringSubmatch(msg)
	if m == nil {
		return nil, false
	}
	return m, true
}

func isGoodHash(hash string) bool {
	switch hash {
		case "SHA256", "SHA512": { }
		default: {
			return false
		}
	}
	return true
}

func loadTextFile(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func isSameNonce(nonceFromUser string, msg string) bool {
	m, _ := makeMatch(msg, nonceRe)
	nonceFromMessage := strings.TrimSpace(m[1])
	return (nonceFromMessage == nonceFromUser)
}

func isValidPgpClearSignature(r url.Values) bool {
	msg := r["Datum"][0]
	var bsmRe = "-----BEGIN PGP SIGNED MESSAGE-----"
	var bpsRe = "-----BEGIN PGP SIGNATURE-----"
	var epsRe = "-----END PGP SIGNATURE-----"
	var hashRe = "Hash: ([A-z0-9-_]*)"
	m, e := makeMatch(msg, hashRe)
	if m != nil {
		if ! isGoodHash(m[1]) {
			return false
		}
	}
	regexes := []string{bsmRe, bpsRe, epsRe, nonceRe}
	ok := true
	for _, v := range regexes {
		_, e = makeMatch(msg, v)
		ok = ok && e
	}
	return ok
}

func isVerifiedPgpClearSignature(r url.Values, user *User, pubkey string) bool {
	if ! isValidPgpClearSignature(r) {
		return false
	}
	nonce := user.Nonce
	user.Nonce = "" // only use once
	datum := r["Datum"][0]
	if ! isSameNonce(nonce, datum) {
		return false
	}
	name := r["User"][0]
	if name == "" {
		return false
	}
	user.Name = name
	verifiedPlainText, err := helper.VerifyCleartextMessageArmored(pubkey, datum, crypto.GetUnixTime())
	if err != nil {
		return false
	}
	if verifiedPlainText != nonce {
		return false
	}
	if err == nil {
		user.Name = name
		user.authorise()
		return true
	}
	return false
}
