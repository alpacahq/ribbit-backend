package mock

import (
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/gin-gonic/gin"
)

// Auth mock
type Auth struct {
	UserFn func(*gin.Context) *model.AuthUser
}

// User mock
func (a *Auth) User(c *gin.Context) *model.AuthUser {
	return a.UserFn(c)
}
