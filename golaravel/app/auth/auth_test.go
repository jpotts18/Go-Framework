package auth

import (
        "testing"
)

func TestHash(t *testing.T) {
        password := "mysecretpassword"
        hash, err := Hash(password)
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }
        if hash == "" {
                t.Error("Expected non-empty hash")
        }
        if hash == password {
                t.Error("Hash should not equal plain password")
        }
}

func TestCheck(t *testing.T) {
        password := "mysecretpassword"
        hash, _ := Hash(password)

        if !Check(password, hash) {
                t.Error("Expected password to match hash")
        }

        if Check("wrongpassword", hash) {
                t.Error("Expected wrong password to not match hash")
        }
}

func TestNeedsRehash(t *testing.T) {
        password := "mysecretpassword"
        hash, _ := Hash(password)

        if NeedsRehash(hash) {
                t.Error("Expected hash to not need rehash")
        }

        if !NeedsRehash("invalid-hash") {
                t.Error("Expected invalid hash to need rehash")
        }
}

func TestGenerateToken(t *testing.T) {
        token1, err := GenerateToken(32)
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }
        if len(token1) == 0 {
                t.Error("Expected non-empty token")
        }

        token2, _ := GenerateToken(32)
        if token1 == token2 {
                t.Error("Expected different tokens")
        }
}

func TestGenerateRememberToken(t *testing.T) {
        token, err := GenerateRememberToken()
        if err != nil {
                t.Fatalf("Expected no error, got: %v", err)
        }
        if len(token) == 0 {
                t.Error("Expected non-empty token")
        }
}

func TestSecureCompare(t *testing.T) {
        if !SecureCompare("abc", "abc") {
                t.Error("Expected equal strings to match")
        }
        if SecureCompare("abc", "def") {
                t.Error("Expected different strings to not match")
        }
}

func TestParseBearerToken(t *testing.T) {
        tests := []struct {
                header   string
                expected string
        }{
                {"Bearer mytoken123", "mytoken123"},
                {"Bearer ", ""},
                {"", ""},
                {"Basic xyz", ""},
        }

        for _, tc := range tests {
                result := ParseBearerToken(tc.header)
                if result != tc.expected {
                        t.Errorf("ParseBearerToken(%q) = %q, want %q", tc.header, result, tc.expected)
                }
        }
}

func TestPasswordValidator(t *testing.T) {
        v := NewPasswordValidator().MinLength(8).RequireUppercase().RequireDigit()

        tests := []struct {
                password string
                valid    bool
        }{
                {"Abc12345", true},
                {"abc12345", false},
                {"ABCDEFGH", false},
                {"Abc", false},
        }

        for _, tc := range tests {
                err := v.Validate(tc.password)
                if tc.valid && err != nil {
                        t.Errorf("Expected %q to be valid, got error: %v", tc.password, err)
                }
                if !tc.valid && err == nil {
                        t.Errorf("Expected %q to be invalid", tc.password)
                }
        }
}

func TestTokenGuard(t *testing.T) {
        resolver := func(token string) (interface{}, interface{}, error) {
                if token == "valid-token" {
                        return map[string]string{"name": "John"}, int64(1), nil
                }
                return nil, nil, nil
        }

        guard := NewTokenGuard(resolver)

        if guard.Check() {
                t.Error("Expected guest before setting token")
        }
        if !guard.Guest() {
                t.Error("Expected guest to be true")
        }

        guard.SetToken("valid-token")

        if !guard.Check() {
                t.Error("Expected authenticated after setting token")
        }
        if guard.Guest() {
                t.Error("Expected guest to be false")
        }
        if guard.User() == nil {
                t.Error("Expected user to be set")
        }
        if guard.ID() != int64(1) {
                t.Errorf("Expected ID to be 1, got %v", guard.ID())
        }

        guard.Logout()
        if guard.Check() {
                t.Error("Expected guest after logout")
        }
}

func TestGuardInterface(t *testing.T) {
        resolver := func(token string) (interface{}, interface{}, error) {
                return nil, nil, nil
        }
        var guard Guard = NewTokenGuard(resolver)
        if guard == nil {
                t.Error("TokenGuard should implement Guard interface")
        }
}
