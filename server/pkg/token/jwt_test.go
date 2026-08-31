package token

import (
	"testing"
	"time"
)

func TestAccessTokenRoundTrip(t *testing.T) {
	secret := "test-secret"
	at, err := GenerateAccessToken(42, false, secret, 30*time.Minute)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if at.Token == "" || at.JTI == "" {
		t.Fatalf("token or jti empty")
	}
	claims, err := ParseAccessToken(at.Token, secret)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("user id = %d, want 42", claims.UserID)
	}
	if claims.Role != "user" {
		t.Fatalf("role = %s, want user", claims.Role)
	}
}

func TestParseWrongSecret(t *testing.T) {
	at, _ := GenerateAccessToken(1, true, "a", time.Minute)
	if _, err := ParseAccessToken(at.Token, "b"); err == nil {
		t.Fatalf("expected error parsing with wrong secret")
	}
}

func TestRandomString(t *testing.T) {
	a, _ := RandomString(16)
	b, _ := RandomString(16)
	if a == "" || len(a) != 32 {
		t.Fatalf("unexpected length: %q", a)
	}
	if a == b {
		t.Fatalf("random strings should differ")
	}
}
