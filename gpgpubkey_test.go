package main
/*
import (
	"testing"
	"net/url"
	"encoding/base64"
)

var nonce1 = "-O-Q7SqA_5Y88V4aTK71"
var bsig1 = `LS0tLS1CRUdJTiBQR1AgU0lHTkVEIE1FU1NBR0UtLS0tLQpIYXNoOiBTSE
EyNTYKCi0gLU8tUTdTcUFfNVk4OFY0YVRLNzEKLS0tLS1CRUdJTiBQR1Ag
U0lHTkFUVVJFLS0tLS0KCmlRSktCQUVCQ0FBMEZpRUVmS0RNR2doWVo0VG
9iMVpEeGRqQ05ZM25aYmtGQW1adDJ1RVdIR1poZEcxaGJrQmsKY21WaGJY
UnlZV05yTG01bGRBQUtDUkRGMk1JMWplZGx1YUtzRUFDZ0lldHNvVVFvYX
FvalFLWGZYNm0rNDFraAo2dGlrYU14ZzlmcjNnTjEwcjgycE5SOWFuQzVu
cHRGeEdxQU9pTFBiclNpVnZkWkNWWmFveVRnM1RnN3BVazA5CndRdkZyem
lWRjFZNmlLZmcwTlJyb2huRGVlU252K2gzUFROUDFJdFpYQjk3QjlUZUR5
VEZQZ1FaSlVBUldBSkMKR1hWQmlBWXczcWF5TVgvQVMwTkVYbWxvZmVvY0
tHM1E3NDJzMy9OdC9aYTlrUmdZRkc5eDh0dmd2V3JRTDZ2dwpVTk9hQnNC
SkprbVVRa2lyaWdtQzgxVG0vamFEcysyeUVQQXlBQ0orbE1xL1g0cVhhZX
ZKdmhQczZNdWE4cGoxCmlueTZHVlppTy9hTy9WZjB0Wm9KYkVsVnFBK0hs
TGlGZXRSbDNPdGtkbkJTRG02N2pVbVhReS9YWEtJdytCa3YKdDZwNjA5Mm
R2UFhnYWJ3TW5pNHBaSXBUREsyWDRWc1RyeC9JVFdxaE9XZDlCZjQrNnN4
M1ZtRHdFZG1sNE9xOApUSVY1dFJ1UGpWbXZwckdQRTh3T2laQ3B5YjRYY2
t3OE5xKzE3d2dzVGJ0N25UZk0wM2s3OUt4UlBnWWdNT3RlCkdHWlBMTHh2
VExmbkJ2dFdVUitsUnhVNE5rU0FibStNVVkvc09sKzY0eGE1Q3V4SXl4S2
VoWlJnaU1NZGtmMVAKaDAvaWloZXNCdFhEdi91cjh1ZFAxdUV2T0VUWUVp
QXdad3k3dmhUbDArcCtsZkJpd25Vb20ycm5KSk1rOWc0RQpkN0dZdUtSKz
d3MGJDWUF2YWEzWlM2SVZUTThEZGltWWZoQW1oS3h1a2RjZEs2cUcvMjd5
ZU5BTTJCbkpmWE9NClRNV1lnV2dFVU5TclpYWUdrQT09Cj1CdDQvCi0tLS
0tRU5EIFBHUCBTSUdOQVRVUkUtLS0tLQo=`

var nonce2 = "-O-Q7SqA_5Y88V4aTK70"
var bsig2 = `LS0tLS1CRUdJTiBQR1AgU0lHTkVEIE1FU1NBR0UtLS0tLQpIYXNoOiBTSE
EyNTYKCi0gLU8tUTdTcUFfNVk4OFY0YVRLNzAKLS0tLS1CRUdJTiBQR1Ag
U0lHTkFUVVJFLS0tLS0KCmlRSktCQUVCQ0FBMEZpRUVmS0RNR2doWVo0VG
9iMVpEeGRqQ05ZM25aYmtGQW1adDJ5b1dIR1poZEcxaGJrQmsKY21WaGJY
UnlZV05yTG01bGRBQUtDUkRGMk1JMWplZGx1V0dMRUFDZmhLL2cyd3huMm
RDd0dRWmovMlJOUEtUWQpxdXIzMFIrOUQ1TWFsbDIwcElNWDl6Sytza3Jm
bEJ3OWpjbkZON1p5TUZtdXIydEt3ZVpyS09UTUZPck9PR0NUCkp1LytTNV
NCWTFZVWt0RWxnaHpmblZxcTBVRlFHcmFhYTdqUXpwbXdUWHJlbEVEM21G
WHdWeHlQUTZwL2FkMi8KSmpaaTJyd1JVamhPbGhIWW1CL080aytKbWQ5M1
JHQWpyZ1N4ZExYeDdXMitKOHl5b0FGNm9WTWhBNHZycHpnWgpPTkJaN1pE
NnZaQUNUWDA0bkdoYU11QmQ0SU9vZEJyWlRoanF4SnFnZ0JWRk9VY0JkS3
c0dWVKaFd2bGpLNEF2CnZsNTZhazdGb0tzMUxXWWRreElLYUJSbHVZY0U3
R3ZlNis2eTk1T1lEUWpDVHZ0NUZ5dkViczRwNXFDeGpNSEgKcjZrdWtJdm
5KdG00NTRITkdlVEs0bU9ZVktCVjhPeVlxbTd3R00rNG9hZUhQMzE2MlhN
OTRVZytrbVI4N3BBTQpUL0UzeXJCSjFPNklJZ0Z6U3hDSWhqeUJxSmJzMT
lMTGxWV2hjYnZkbEgxVkh4ejhsODFTeVFtNUhMU01DeFVjCnlpZTVEd2tO
SkI0OTdCQlB4V0ZTRmRQUFc4Z21leXJJeFJDNW5pRm5DQ0VaVVR0VjhQcW
p6L2RNcmd3NnM0NGcKZkxLWU1wMkVjdzJRNEZTVkxmR3ZlV2hHc01QaVJL
OWsvRWYwTXpzUjZsRDhseFFURlZCY3FmcmVqaGhSc2ZiVApCRnhmUENzRy
9Zb2x2RU1neUFMZWVXVFFKNmxWT0Q0L2lDbHU1c1VjNHVPeUNaMktmSTBj
QXltak9jUVJSVmxQCmZ2U1Z3YTRwc2xxZXZJMUxZQT09Cj1Gd2RZCi0tLS
0tRU5EIFBHUCBTSUdOQVRVUkUtLS0tLQo=`

var bpubkey = `LS0tLS1CRUdJTiBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0tCgptUUlO
QkdQVUROd0JFQUM0SWVmQmY2MktXTTNRRUxsWUt5VE1nb2RRY2xndWlI
SE5Nbm1UTDFsbzZBbWV2Rm9mCkdoSWJZSXlyQXJNcWxKdEhlNTVjWjlM
ZlJEN2wzTFl0TVpsU1Q2WnI0NjI5eThEU0d0RHQvU2NZQk9KNk9memYK
Zkk5YncvRno3U01yVU41TTNyOWhNVFpPR1kyOE5KcUNncDlLam5NUFZq
bE91NmhvNDlBMDNkQXg3L29xS2NOdwpwaFROa09pUXNqMXJ1dk8wY2lO
RVZoYklXdFdPNjNZZExIbkkwWFF1dmhBSVN2TnVLcU5TcWtMUWY3Qzkx
S2lpCmFGWDZNeVRhdFFoN0hQdGxJalIzbnBJVDhCaEd0anBXTDl5aXZD
ZkRmbExiZ3RRUi9NT1J6cTYrRHE2dmt2czgKK0ZGa2g2SEk0ZnJvaTNT
cDdmQkFnWlhZb2xKTmE1dlFTaUZQSVNXZzVabDllanNjcWZGTTBiaWcz
UTByLzg5RwpkVnZDWnRzUXowb05RSVhVUlNNMFZVQzRQbzRFdEdUdjFX
WkJFU0llNWxtdDZhMnJHbkJYd2tTRXpoQTZsOVNJClJ6T21FQUxuWGdj
NTJXVjNlTjR3TVhTVnR0dlJRb2dweVFOT0RRZ2VWaGI0RGRSeFk3SE1R
U0FLK1JWS0NVV2YKNWx0K083NkR6ZFZwUXQ0dzMwSjNmRk5pREFKb0xa
RlR6am50clVkNGxkcGhYeWtMa2FUYTZlcVFmeEJRb2ZEWApoZjhjYzNH
Y0FHZjZYWHdBTWlaa21EMHdKZzFYWDhSbCsycFFaNXZxc2NJZG5kbzdX
MmljNzlacjQzNEhRbVpSCnFQd1k0L0tWVnBaeXhuRHJsbUNCWURYaVJh
WmhrRSt4aldDSmZYRkdicm1vRHhycjNORTZrSTdMWXdBUkFRQUIKdENk
QlpHRnRJRkpwWTJoaGNtUnpiMjRnUEdaaGRHMWhia0JrY21WaGJYUnlZ
V05yTG01bGRENkpBbGNFRXdFSQpBRUVDR3dNRkN3a0lCd0lDSWdJR0ZR
b0pDQXNDQkJZQ0F3RUNIZ2NDRjRBV0lRUjhvTXdhQ0ZobmhPaHZWa1BG
CjJNSTFqZWRsdVFVQ1pqME5DUVVKQkVvenJRQUtDUkRGMk1JMWplZGx1
VUo5RC93Tmh4N1ZubFZuQjlaRTFtRkIKZ2pLU1lTa2l0UDMvRG12aCtF
b2h4bzVzTHphNFF3THBLMUhBK0pGRzVyTmdsdElVVUw4TWRCNE05N2kr
YUlaYQphK2tsWGJBWWt5b0Q3N1U4MHEwdU1rK1NUVEc1Qmg1TVgzMDZv
ZzNNY3ZnK3JucHhSTXVVSTBZRGhrdm1MMzNvCjNIKytXYmQ0eWdveUhq
TWRrdE1ocjl5S0NOcnUxV0V2ZEo3Q3ZuQSsvdzRJT01xQ0RNT2JoSEp4
Y0YxcUxJSmkKOHYzWG5uUitpU2xtTkRxMjl6c0xXQ1Znb3plT3QybUJE
S3FoYjNhQWVUalc4b04rclV5a1ZhTkRzU2NkK2NKUwpCaVdKOU54RXZn
L0lDcWdub2hzQzd4Yys4ZnYvOG9vU2dQVDdUZ09aVGNuOGtBdnA4aDVa
NFJlb0RJOXFlK2R0CjlFOS8zT1VSdXB6aHFXSlhnOHA5cVBBbURodllC
d2Y4YjhGNWkxeld2WDNpTTlsSlFkYlNkeTFaN3lVVVRiMjIKRDQxdGVm
YnFydkhNK2trRk9hcnZXK3BZKzJlOUdBRTlWVUlnSEtlZWI0OEZHdUxU
cWVuR1FyTnJpdnIwUUZkMgoxd0NCWWJOZVZXMzk2eXY3YlkzajZEblZl
cjJYNk1GNmtvWTlaS1J6akwzT3V1Z1dsWUlFVXBpbktGdkEzMW9sCnFl
R3htUFQ0TWI1SGh1TnJDdWx3YkMza1dTL2hTYkcxZzFDUjdKY2xUTkk3
cFRFcG9oZ3AvMHNSMlFETDZ4R3gKOExjTEJWSFJTbjQ0bG9nMmY4Y3pu
Wmc3K0tyY1IvaHpQVmZubUFrekhhQmJDenVrNkdGVDFOcDBqSHJQcTNk
ZgpOU2YyR01nNmhJWGR3aUtIeWE0Z1FWNFNmTGtDRFFSajFBemNBUkFB
OU5TL2w4cHNJS3BuT3VpU0ZOcUZBZ0RECmlEZndGMWhZZXZ1d1haaEVY
N3Q0NE9JT29MWkl2aTFuUUVWcnNVanJGcTQxR0Q1SzlsMllrNXMxbTRH
TXdNOWIKeW9rS3pEb3YyY1NndXB2c1ZuTEZmUENGTzN5UmJNS2h2NDlp
NkVDR3RxZ3lDRGNiUVdVTkJ4enJZTXRPak56egpicld4Nmw2bVVySmwr
eXB5SkJnV1piMmJmNHJsUlkyTGYwU1pNdUVWbW01clg5MDVYMUJoQ296
NjdFNlo1R2hCCjZhNWt6aDV5SDhNY0c4NjhWbzQ3V0FlSlZ5U3oxVlB6
MVhjSkkwd3Z6S3pENGE4MVhHTUs3OGE5Q0MrcFJVTGkKVm9xMzJUWjN3
LzREeFpGOTlPVEt1Zi9YVVNaeEc3R3orMXYxRXZ4cGVyWTh1TUNiWkEx
ZktNZkh3dS9YYTYzcQpubFY5dTU2MjhFR2x1RUl1cEozQ1ZQejlMUGxK
RzZoTno2cWhmY1kwWDNQRS9OVVdmbFZraXorVU9meGVBWkNLCnhDdGVV
c0w2d3NhOERZc3RJK1BJUStkNnNlVmlUMStpU1ZMRHhBcURVb010SG5h
VG1wZ1Z6T0lUN2RERzFEMUsKdnRGQ3hxZzc1ZkVrNE9sU1JRb2Q3UjNL
ZGFXUFVDRTlyZjZBU0l5djlWVXJXemRkSEd6eGltMmpLNTlNTDY5egpy
aDBJeWF2bnl4R3paUGZDMnNFcjkzb3E3ZUtkajlVUHFOalZXL2s3M3Mr
UHZmS3JHbTdUN0swQ2Q5WnRvbzZoCnlRZEp2VS9IS2hTUWJUZzlDTGpq
cWRqbE84TmVnYUYzZ056b3VRRnBSSzhwTGZJN2Rycmcxd3hnalpoZG1G
azMKem5wU3NWV3h5K08yTmpRazdDMEFFUUVBQVlrQ1BBUVlBUWdBSmdJ
YkRCWWhCSHlnekJvSVdHZUU2RzlXUThYWQp3aldONTJXNUJRSm1QUTJK
QlFrRVNqUXRBQW9KRU1YWXdqV041Mlc1c09zUC8zdjBvOURMSE00TGF5
ZExjRDF0ClRvVFJRY0svZjZTczFJeDh2b3VDSm9HeFNvL2I1SG9LczJx
K0tIYTl4N1FqRlNuWk5Gdi80VjBNdHRDdkVnTngKeXAxdkp2eTVCeG1y
SGRkbTY5TVJRbitMS3BqS3Y5WFpHbmprbXdaS1F5TlhBeE5NOEhWa21r
ZWJDckpUQTNFdwpNeG5UNE50dDZpSnRlZDRKQlFiREo4SS91UkxPT2NN
SytxWEM5TVNkNURCMlorLzRXOFVSa29hL09MczdKUlNUCmpwSVdiZ0xZ
SFRySUQ2QWRmR2dpZkJsYzk5dTg3SEFKdUFyZWk3SzFoV3ZMK1VtaDJD
TGRHOE13NkRpREZ2UXgKTnh2Mk9uVzFWc2YvWHRqWHFpNWtzd29NOHdw
aVg2ajNsaTkyRmRZU2hIUm9oS2lPSzV4QW9oK0lia2hzMVJvMgpxdUJC
N3hUWTArVDI2U3ZLYUczUFY5MG5ydjJWV3ZMcC9oNnRUUFBsTHMvdC9u
MjJzSXljOW5SVXU5ZjUveEtJCktEd0lWUkFMcVc5L203T3g4aW1KR1Ju
dDIwMkpxTFpkRnc0T0ZuK1lsMm82djhVNnk2Qm9sUmJMaXJ1R3hyRHoK
RGR6TEEyTHJrZ0JvTjNPWWZrOXRFRUIyTVZaVUxrb296ZW1XeTRkNWl3
U1FhdVZsWkQybURUOUJESUVkamk0UwpLY09NRVg2VnZCeTBRa0lQWHlG
d201aFV1STNFb1J2Qk14L1N1Q2d4Tzd4OEIrakxJbWZBZnRDbE1tRUNm
Uzh1CkxkK2xEaXlsakUzNGg3MGZQTSthYXdzb09XL09QRXNZaXRhdWVS
TUU4SnEvdWM5TjVRK2hjZVZNcW4wbVppOTYKUHVtYUtQZ2xWYzFQZ09S
QWw1bm9YZExOCj05eGtjCi0tLS0tRU5EIFBHUCBQVUJMSUMgS0VZIEJM
T0NLLS0tLS0K`

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
	dec1, err := base64.StdEncoding.DecodeString(bsig1)
	if err != nil {
		panic(err.Error())
	}
	dec2, err := base64.StdEncoding.DecodeString(bsig2)
	if err != nil {
		panic(err.Error())
	}
	if ! isSameNonce(nonce2, string(dec1)) {
		t.Logf("-O-Q7SqA_5Y88V4aTK70: false (ok)")
	} else {
		t.Fatalf("Should be false")
	}
	if isSameNonce(nonce2, string(dec2)) {
		t.Logf("-O-Q7SqA_5Y88V4aTK70: true (ok)")
	} else {
		t.Fatalf("Should be true")
	}
}

func TestValidPgpClearSignature(t *testing.T) {
	dec1, err := base64.StdEncoding.DecodeString(bsig1)
	if err != nil {
		panic(err.Error())
	}
	dec2, err := base64.StdEncoding.DecodeString(bsig2)
	if err != nil {
		panic(err.Error())
	}
	v1 := url.Values{}
	v1.Add("User", "fatman@dreamtrack.net")
	v1.Add("Datum", string(dec1))
	v2 := url.Values{}
	v2.Add("User", "fatman@dreamtrack.net")
	v2.Add("Datum", string(dec2))
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
	dec1, err := base64.StdEncoding.DecodeString(bsig1)
	if err != nil {
		panic(err.Error())
	}
	dec2, err := base64.StdEncoding.DecodeString(bsig2)
	if err != nil {
		panic(err.Error())
	}
	decpubkey, err := base64.StdEncoding.DecodeString(bpubkey)
	if err != nil {
		panic(err.Error())
	}
	v1 := url.Values{}
	v1.Add("User", "fatman@dreamtrack.net")
	v1.Add("Datum", string(dec1))
	v2 := url.Values{}
	v2.Add("User", "fatman@dreamtrack.net")
	v2.Add("Datum", string(dec2))
	u1 := User{ Name: "fatman@dreamtrack.net", Nonce: nonce1}
	u2 := User{ Name: "fatman@dreamtrack.net", Nonce: nonce2}
	if isVerifiedPgpClearSignature(v1, &u1, string(decpubkey)) {
		t.Logf("sig1 is good")
	} else {
		t.Fatalf("sig1 is bad")
	}
	if isVerifiedPgpClearSignature(v2, &u2, string(decpubkey)) {
		t.Logf("sig2 is good")
	} else {
		t.Fatalf("sig2 is bad")
	}
}
*/
