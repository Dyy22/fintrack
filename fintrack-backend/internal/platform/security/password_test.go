package security

import "testing"

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "securepassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty password hash")
	}
	if hash == password {
		t.Fatal("password hash must not equal plain password")
	}
	if !CheckPassword(hash, password) {
		t.Fatal("expected password check to pass for matching password")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Fatal("expected password check to fail for wrong password")
	}
}
