package util_test

import (
	"testing"

	"github.com/sejamuchhal/taskhub/auth/util"
	"golang.org/x/crypto/bcrypt"
)

func TestValidPassword(t *testing.T) {
	password := "test#2dff"
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	err = util.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Password should be valid, but check failed: %v", err)
	}
}

func TestInvalidPassword(t *testing.T) {
	password := "test#2dff"
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	invalidPassword := "wrongPsW"
	err = util.CheckPasswordHash(invalidPassword, hashedPassword)
	if err == nil {
		t.Error("Invalid password did not fail")
	}
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("Expected error %v, but got %v", bcrypt.ErrMismatchedHashAndPassword, err)
	}
}
