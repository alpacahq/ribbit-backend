package mock

import "github.com/alpacahq/ribbit-backend/model"

// JWT mock
type JWT struct {
	GenerateTokenFn func(*model.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *model.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
