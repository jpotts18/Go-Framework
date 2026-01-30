package auth

import (
        "crypto/rand"
        "crypto/subtle"
        "encoding/base64"
        "fmt"
        "strings"

        "golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
        bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                return "", err
        }
        return string(bytes), nil
}

func Check(password, hash string) bool {
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        return err == nil
}

func NeedsRehash(hash string) bool {
        cost, err := bcrypt.Cost([]byte(hash))
        if err != nil {
                return true
        }
        return cost < bcrypt.DefaultCost
}

func GenerateToken(length int) (string, error) {
        bytes := make([]byte, length)
        if _, err := rand.Read(bytes); err != nil {
                return "", err
        }
        return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateRememberToken() (string, error) {
        return GenerateToken(32)
}

func SecureCompare(a, b string) bool {
        return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type Guard interface {
        Check() bool
        Guest() bool
        User() interface{}
        ID() interface{}
        Login(user interface{}) error
        Logout() error
}

type TokenGuard struct {
        token    string
        user     interface{}
        userID   interface{}
        resolver func(token string) (interface{}, interface{}, error)
}

func NewTokenGuard(resolver func(token string) (interface{}, interface{}, error)) *TokenGuard {
        return &TokenGuard{
                resolver: resolver,
        }
}

func (g *TokenGuard) SetToken(token string) {
        g.token = token
        if token != "" && g.resolver != nil {
                user, id, err := g.resolver(token)
                if err == nil {
                        g.user = user
                        g.userID = id
                }
        }
}

func (g *TokenGuard) Check() bool {
        return g.user != nil
}

func (g *TokenGuard) Guest() bool {
        return g.user == nil
}

func (g *TokenGuard) User() interface{} {
        return g.user
}

func (g *TokenGuard) ID() interface{} {
        return g.userID
}


func (g *TokenGuard) Login(user interface{}) error {
        g.user = user
        return nil
}

func (g *TokenGuard) Logout() error {
        g.user = nil
        g.userID = nil
        g.token = ""
        return nil
}

func ParseBearerToken(header string) string {
        if header == "" {
                return ""
        }
        if !strings.HasPrefix(header, "Bearer ") {
                return ""
        }
        return strings.TrimPrefix(header, "Bearer ")
}

func MakePasswordRules() []string {
        return []string{"required", "min:8", "confirmed"}
}

type PasswordValidator struct {
        minLength    int
        requireUpper bool
        requireLower bool
        requireDigit bool
}

func NewPasswordValidator() *PasswordValidator {
        return &PasswordValidator{
                minLength:    8,
                requireUpper: false,
                requireLower: false,
                requireDigit: false,
        }
}

func (v *PasswordValidator) MinLength(n int) *PasswordValidator {
        v.minLength = n
        return v
}

func (v *PasswordValidator) RequireUppercase() *PasswordValidator {
        v.requireUpper = true
        return v
}

func (v *PasswordValidator) RequireLowercase() *PasswordValidator {
        v.requireLower = true
        return v
}

func (v *PasswordValidator) RequireDigit() *PasswordValidator {
        v.requireDigit = true
        return v
}

func (v *PasswordValidator) Validate(password string) error {
        if len(password) < v.minLength {
                return fmt.Errorf("password must be at least %d characters", v.minLength)
        }

        if v.requireUpper {
                hasUpper := false
                for _, c := range password {
                        if c >= 'A' && c <= 'Z' {
                                hasUpper = true
                                break
                        }
                }
                if !hasUpper {
                        return fmt.Errorf("password must contain at least one uppercase letter")
                }
        }

        if v.requireLower {
                hasLower := false
                for _, c := range password {
                        if c >= 'a' && c <= 'z' {
                                hasLower = true
                                break
                        }
                }
                if !hasLower {
                        return fmt.Errorf("password must contain at least one lowercase letter")
                }
        }

        if v.requireDigit {
                hasDigit := false
                for _, c := range password {
                        if c >= '0' && c <= '9' {
                                hasDigit = true
                                break
                        }
                }
                if !hasDigit {
                        return fmt.Errorf("password must contain at least one digit")
                }
        }

        return nil
}
