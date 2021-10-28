package secret

import (
	"golang.org/x/crypto/bcrypt"
)

// New returns a password object
func New() *Password {
	return &Password{}
}

// Password is our secret service implementation
type Password struct{}

// HashPassword hashes the password using bcrypt
func (p *Password) HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func (p *Password) HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// HashRandomPassword creates a random password for passwordless mobile signup
func (p *Password) HashRandomPassword() (string, error) {
	randomPassword, err := GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	r := p.HashPassword(randomPassword)
	return r, nil
}
