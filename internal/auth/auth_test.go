package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	samplePassword := "password123"
	hash, err := HashPassword(samplePassword)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}
	err = CheckPasswordHash("wrongPassword", hash)
	if err == nil {
		t.Fatalf("password check should have failed")
	}
	err = CheckPasswordHash("password123", hash)
	if err != nil {
		t.Fatalf("password check shouldn't have failed")
	}
}

func TestJWT(t *testing.T) {
	uuid := uuid.New()
	jwtSecret := "password"
	jwt, err := MakeJWT(uuid, jwtSecret, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}
	t.Logf("\njwt: %v", jwt)
	result, err := ValidateJWT(jwt, jwtSecret)
	if err != nil {
		t.Fatalf("error validating token: %v", err)
	}
	if result != uuid {
		t.Fatalf("%s != %s", uuid, result)
	}
	time.Sleep(300 * time.Millisecond)
	_, err = ValidateJWT(jwt, jwtSecret)
	if err == nil {
		t.Fatalf("token should have been expired")
	}
	_, err = ValidateJWT(jwt, "wrongSecret")
	if err == nil {
		t.Fatalf("token should not have been verified")
	}
}

func TestGetBearerToken(t *testing.T) {
	token := "123"
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	result, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("could not get bearer token: %v", err)
	}
	if result != token {
		t.Fatalf("token isn't correct: %s != %s", result, token)
	}
}
