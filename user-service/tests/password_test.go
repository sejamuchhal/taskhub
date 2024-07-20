package tests

import (
	"testing"

	"github.com/sejamuchhal/task-management/user-service/internal/server"
	"golang.org/x/crypto/bcrypt"
)

func TestValidPassword(t *testing.T) {
	password := "test#2dff"
	hashedPassword, err := server.HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	err = server.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Password should be valid, but check failed: %v", err)
	}
}

func TestInvalidPassword(t *testing.T) {
	password := "test#2dff"
	hashedPassword, err := server.HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	invalidPassword := "wrongPsW"
	err = server.CheckPasswordHash(invalidPassword, hashedPassword)
	if err == nil {
		t.Error("Invalid password did not fail")
	}
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("Expected error %v, but got %v", bcrypt.ErrMismatchedHashAndPassword, err)
	}
}
