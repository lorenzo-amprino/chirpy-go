package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	match, err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !match {
		t.Fatal("Expected password to match hash, but it did not")
	}

	wrongPassword := "wrongpassword"
	match, err = CheckPasswordHash(wrongPassword, hashed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match {
		t.Fatal("Expected password to not match hash, but it did")
	}
}

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hashed == password {
		t.Fatal("Expected hashed password to be different from original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	match, err := CheckPasswordHash(password, hashed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !match {
		t.Fatal("Expected password to match hash, but it did not")
	}
}

func TestCheckPasswordHashWithWrongPassword(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	wrongPassword := "wrongpassword"
	match, err := CheckPasswordHash(wrongPassword, hashed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match {
		t.Fatal("Expected password to not match hash, but it did")
	}
}
