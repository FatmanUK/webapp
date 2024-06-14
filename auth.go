package main

import (
	"net/http"
	"regexp"
)

func makeMatch(msg string, re string) ([]string, bool) {
	r := regexp.MustCompile(re)
	m := r.FindStringSubmatch(msg)
	if m == nil {
		return nil, false
	}
	return m, true
}

func isFilePath(f string) bool {
	// TODO
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
	var line4Re = "(.*?) (.*?)"

	ok := true

	_, e := makeMatch(msg, bsmRe)
	ok = ok && e

	_, e = makeMatch(msg, bpsRe)
	ok = ok && e

	_, e = makeMatch(msg, epsRe)
	ok = ok && e

	m, e := makeMatch(msg, hashRe)
	ok = ok && e

	if m != nil {
		hash := m[1]
		switch hash {
			case "SHA256", "SHA512": { }
			default: {
				ok = false
			}
		}
	}

	m, e = makeMatch(msg, line4Re)
	ok = ok && e

	if m != nil {
		file := m[1]
		if ! (isFilePath(file) || file == "-") {
			ok = false
		}
	}
	return ok
}

func sigIsVerified(name string, datum string, nonce string) bool {
	if ! sameNonce(nonce, datum) {
		return false
	}

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
	// https://stackoverflow.com/questions/33963284/verify-gpg-signature-in-go-openpgp
	return true
}

func isVerifiedPgpClearSignature(r *http.Request, user *User) bool {
	if ! isValidPgpClearSignature(r) {
		return false
	}
	nonce := user.Nonce
	user.Nonce = "" // only use once
	name := r.PostForm["User"][0]
	datum := r.PostForm["Datum"][0]
	if sigIsVerified(name, datum, nonce) {
		user.Name = name
		user.authorise()
		return true
	}
	return false
}
