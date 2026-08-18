package auth

import (
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == password {
		t.Fatal("HashPassword() returned the password in plain text")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if !match {
		t.Error("CheckPasswordHash() = false for the correct password, want true")
	}

	match, err = CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() error = %v", err)
	}
	if match {
		t.Error("CheckPasswordHash() = true for a wrong password, want false")
	}
}

func TestMakeJWTAndValidate(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}
	if gotID != userID {
		t.Errorf("ValidateJWT() = %v, want %v", gotID, userID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "super-secret", -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	gotID, err := ValidateJWT(token, "super-secret")
	if err == nil {
		t.Fatal("ValidateJWT() error = nil for an expired token, want an error")
	}
	if gotID != uuid.Nil {
		t.Errorf("ValidateJWT() = %v for an expired token, want uuid.Nil", gotID)
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "super-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	gotID, err := ValidateJWT(token, "not-the-secret")
	if err == nil {
		t.Fatal("ValidateJWT() error = nil for a token signed with another secret, want an error")
	}
	if gotID != uuid.Nil {
		t.Errorf("ValidateJWT() = %v for a wrongly signed token, want uuid.Nil", gotID)
	}
}

func TestValidateJWTRejectsMalformedToken(t *testing.T) {
	if _, err := ValidateJWT("not.a.jwt", "super-secret"); err == nil {
		t.Fatal("ValidateJWT() error = nil for a malformed token, want an error")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token := MakeRefreshToken()

	if len(token) != 64 {
		t.Errorf("MakeRefreshToken() length = %d, want 64 hex chars (32 bytes)", len(token))
	}
	if _, err := hex.DecodeString(token); err != nil {
		t.Errorf("MakeRefreshToken() = %q, not valid hex: %v", token, err)
	}
	if other := MakeRefreshToken(); other == token {
		t.Error("MakeRefreshToken() returned the same token twice, want random values")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		setHeader bool
		want      string
		wantErr   bool
	}{
		{name: "valid token", header: "Bearer TOKEN_STRING", setHeader: true, want: "TOKEN_STRING"},
		{name: "extra whitespace", header: "Bearer   TOKEN_STRING  ", setHeader: true, want: "TOKEN_STRING"},
		{name: "no header", setHeader: false, wantErr: true},
		{name: "empty header", header: "", setHeader: true, wantErr: true},
		{name: "missing bearer prefix", header: "TOKEN_STRING", setHeader: true, wantErr: true},
		{name: "wrong scheme", header: "Basic TOKEN_STRING", setHeader: true, wantErr: true},
		{name: "bearer without token", header: "Bearer   ", setHeader: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.setHeader {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetBearerToken(headers)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetBearerToken() error = nil, want an error (got %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetBearerToken() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		setHeader bool
		want      string
		wantErr   bool
	}{
		{name: "valid key", header: "ApiKey THE_KEY_HERE", setHeader: true, want: "THE_KEY_HERE"},
		{name: "extra whitespace", header: "ApiKey   THE_KEY_HERE  ", setHeader: true, want: "THE_KEY_HERE"},
		{name: "no header", setHeader: false, wantErr: true},
		{name: "empty header", header: "", setHeader: true, wantErr: true},
		{name: "missing apikey prefix", header: "THE_KEY_HERE", setHeader: true, wantErr: true},
		{name: "wrong scheme", header: "Bearer THE_KEY_HERE", setHeader: true, wantErr: true},
		{name: "wrong case", header: "apikey THE_KEY_HERE", setHeader: true, wantErr: true},
		{name: "apikey without key", header: "ApiKey   ", setHeader: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.setHeader {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetAPIKey(headers)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetAPIKey() error = nil, want an error (got %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetAPIKey() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetAPIKey() = %q, want %q", got, tt.want)
			}
		})
	}
}
