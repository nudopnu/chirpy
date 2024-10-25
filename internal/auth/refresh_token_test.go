package auth

import (
	"log"
	"testing"
)

func TestMake(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		log.Fatalf("error generating token: %v", err)
	}
	if len(token) != 64 {
		log.Fatalf("token '%s' has the worng length! %d != 64", token, len(token))
	}
}
