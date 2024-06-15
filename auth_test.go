package main

import (
	"testing"
	"net/url"
)

var nonce1 = "-O-Q7SqA_5Y88V4aTK71"
var sig1 = `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

- -O-Q7SqA_5Y88V4aTK71
-----BEGIN PGP SIGNATURE-----

iQJKBAEBCAA0FiEEfKDMGghYZ4Tob1ZDxdjCNY3nZbkFAmZt2uEWHGZhdG1hbkBk
cmVhbXRyYWNrLm5ldAAKCRDF2MI1jedluaKsEACgIetsoUQoaqojQKXfX6m+41kh
6tikaMxg9fr3gN10r82pNR9anC5nptFxGqAOiLPbrSiVvdZCVZaoyTg3Tg7pUk09
wQvFrziVF1Y6iKfg0NRrohnDeeSnv+h3PTNP1ItZXB97B9TeDyTFPgQZJUARWAJC
GXVBiAYw3qayMX/AS0NEXmlofeocKG3Q742s3/Nt/Za9kRgYFG9x8tvgvWrQL6vw
UNOaBsBJJkmUQkirigmC81Tm/jaDs+2yEPAyACJ+lMq/X4qXaevJvhPs6Mua8pj1
iny6GVZiO/aO/Vf0tZoJbElVqA+HlLiFetRl3OtkdnBSDm67jUmXQy/XXKIw+Bkv
t6p6092dvPXgabwMni4pZIpTDK2X4VsTrx/ITWqhOWd9Bf4+6sx3VmDwEdml4Oq8
TIV5tRuPjVmvprGPE8wOiZCpyb4Xckw8Nq+17wgsTbt7nTfM03k79KxRPgYgMOte
GGZPLLxvTLfnBvtWUR+lRxU4NkSAbm+MUY/sOl+64xa5CuxIyxKehZRgiMMdkf1P
h0/iihesBtXDv/ur8udP1uEvOETYEiAwZwy7vhTl0+p+lfBiwnUom2rnJJMk9g4E
d7GYuKR+7w0bCYAvaa3ZS6IVTM8DdimYfhAmhKxukdcdK6qG/27yeNAM2BnJfXOM
TMWYgWgEUNSrZXYGkA==
=Bt4/
-----END PGP SIGNATURE-----`

var nonce2 = "-O-Q7SqA_5Y88V4aTK70"
var sig2 = `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

- -O-Q7SqA_5Y88V4aTK70
-----BEGIN PGP SIGNATURE-----

iQJKBAEBCAA0FiEEfKDMGghYZ4Tob1ZDxdjCNY3nZbkFAmZt2yoWHGZhdG1hbkBk
cmVhbXRyYWNrLm5ldAAKCRDF2MI1jedluWGLEACfhK/g2wxn2dCwGQZj/2RNPKTY
qur30R+9D5Mall20pIMX9zK+skrflBw9jcnFN7ZyMFmur2tKweZrKOTMFOrOOGCT
Ju/+S5SBY1YUktElghzfnVqq0UFQGraaa7jQzpmwTXrelED3mFXwVxyPQ6p/ad2/
JjZi2rwRUjhOlhHYmB/O4k+Jmd93RGAjrgSxdLXx7W2+J8yyoAF6oVMhA4vrpzgZ
ONBZ7ZD6vZACTX04nGhaMuBd4IOodBrZThjqxJqggBVFOUcBdKw4ueJhWvljK4Av
vl56ak7FoKs1LWYdkxIKaBRluYcE7Gve6+6y95OYDQjCTvt5FyvEbs4p5qCxjMHH
r6kukIvnJtm454HNGeTK4mOYVKBV8OyYqm7wGM+4oaeHP3162XM94Ug+kmR87pAM
T/E3yrBJ1O6IIgFzSxCIhjyBqJbs19LLlVWhcbvdlH1VHxz8l81SyQm5HLSMCxUc
yie5DwkNJB497BBPxWFSFdPPW8gmeyrIxRC5niFnCCEZUTtV8Pqjz/dMrgw6s44g
fLKYMp2Ecw2Q4FSVLfGveWhGsMPiRK9k/Ef0MzsR6lD8lxQTFVBcqfrejhhRsfbT
BFxfPCsG/YolvEMgyALeeWTQJ6lVOD4/iClu5sUc4uOyCZ2KfI0cAymjOcQRRVlP
fvSVwa4pslqevI1LYA==
=FwdY
-----END PGP SIGNATURE-----`

var pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
PumaKPglVc1PgORAl5noXdLN
=9xkc
-----END PGP PUBLIC KEY BLOCK-----`

func TestGoodHash(t *testing.T) {
	if isGoodHash("SHA256") {
		t.Logf("SHA256 ok")
	} else {
		t.Fatalf("Should be true")
	}
	if isGoodHash("SHA512") {
		t.Logf("SHA512 ok")
	} else {
		t.Fatalf("Should be true")
	}
}

func TestSameNonce(t *testing.T) {
	if ! isSameNonce(nonce2, sig1) {
		t.Logf("-O-Q7SqA_5Y88V4aTK70: false (ok)")
	} else {
		t.Fatalf("Should be false")
	}
	if isSameNonce(nonce2, sig2) {
		t.Logf("-O-Q7SqA_5Y88V4aTK70: true (ok)")
	} else {
		t.Fatalf("Should be true")
	}
}

func TestValidPgpClearSignature(t *testing.T) {
	v1 := url.Values{}
	v1.Add("User", "fatman@dreamtrack.net")
	v1.Add("Datum", sig1)
	v2 := url.Values{}
	v2.Add("User", "fatman@dreamtrack.net")
	v2.Add("Datum", sig2)
	if isValidPgpClearSignature(v1) {
		t.Logf("sig1 is valid")
	} else {
		t.Fatalf("sig1 is not valid")
	}
	if isValidPgpClearSignature(v2) {
		t.Logf("sig2 is valid")
	} else {
		t.Fatalf("sig2 is not valid")
	}
}

func TestVerifiedPgpClearSignature(t *testing.T) {
	// type Values map[string][]string
	v1 := url.Values{}
	v1.Add("User", "fatman@dreamtrack.net")
	v1.Add("Datum", sig1)
	v2 := url.Values{}
	v2.Add("User", "fatman@dreamtrack.net")
	v2.Add("Datum", sig2)
	u1 := User{ Name: "fatman@dreamtrack.net", Nonce: nonce1}
	u2 := User{ Name: "fatman@dreamtrack.net", Nonce: nonce2}
	if isVerifiedPgpClearSignature(v1, &u1, pubkey) {
		t.Logf("sig1 is good")
	} else {
		t.Fatalf("sig1 is bad")
	}
	if isVerifiedPgpClearSignature(v2, &u2, pubkey) {
		t.Logf("sig2 is good")
	} else {
		t.Fatalf("sig2 is bad")
	}
}
