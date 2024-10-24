package auth

import "testing"

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
