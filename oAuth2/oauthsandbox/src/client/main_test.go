package main

import "testing"
import "log"

func TestMakeCodeChallenge(t *testing.T) {
	codeVerifier := "code-challange43128unreserved-._~nge43128dX"
	challenge := makeCodeChallenge(codeVerifier)
	log.Println(codeVerifier)
	log.Println(challenge)
	// SHA256 8B0C93166CF2B27E00C8D52FC3276DF06338893BAEBA6975B17167DC2E939380
	// return "OEIwQzkzMTY2Q0YyQjI3RTAwQzhENTJGQzMyNzZERjA2MzM4ODkzQkFFQkE2OTc1QjE3MTY3REMyRTkzOTM4MA=="
}