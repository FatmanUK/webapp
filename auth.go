package main

import (
	"net/http"
	"regexp"
	"fmt"
	"strings"
//	"golang.org/x/crypto/openpgp"
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

func isSameNonce(nonceFromUser string, msg string) bool {
	m, _ := makeMatch(msg, nonceRe)
	nonceFromMessage := strings.TrimSpace(m[1])
	return (nonceFromMessage == nonceFromUser)
}

func isValidPgpClearSignature(r *http.Request) bool {
	msg := r.PostForm["Datum"][0]
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

func isVerifiedPgpClearSignature(r *http.Request, user *User) bool {
	if ! isValidPgpClearSignature(r) {
		return false
	}
	nonce := user.Nonce
	user.Nonce = "" // only use once
	datum := r.PostForm["Datum"][0]
	if ! isSameNonce(nonce, datum) {
		return false
	}
	//fmt.Println("Good nonce.")
	name := r.PostForm["User"][0]
	user.Name = name

//	public_key_armored := "keys/" + name + ".asc"

// https://stackoverflow.com/questions/33963284/verify-gpg-signature-in-go-openpgp
/*
$ <<<"123456789" gpg --clear-sign | gpg --verify
gpg: Signature made Wed 12 Jun 2024 13:07:21 BST
gpg:                using RSA key B2114...0310E8
gpg: Good signature from "Adam Richardson <adam.richardson@>" [ultimate]

Modified hashed string a little...
$ <myfile gpg --verify
gpg: Signature made Wed 12 Jun 2024 13:07:31 BST
gpg:                using RSA key B2114...0310E8
gpg: BAD signature from "Adam Richardson <adam.richardson@>" [ultimate]
*/
//	if isVerified(name, datum, nonce)
	{
		user.Name = name
		user.authorise()
		return true
	}
	return false
}
