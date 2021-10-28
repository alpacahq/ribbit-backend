package mock

import (
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/gin-gonic/gin"
)

// RBAC Mock
type RBAC struct {
	EnforceRoleFn   func(*gin.Context, model.AccessRole) bool
	EnforceUserFn   func(*gin.Context, int) bool
	AccountCreateFn func(*gin.Context, int) bool
	IsLowerRoleFn   func(*gin.Context, model.AccessRole) bool
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c *gin.Context, role model.AccessRole) bool {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c *gin.Context, id int) bool {
	return a.EnforceUserFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c *gin.Context, roleID int) bool {
	return a.AccountCreateFn(c, roleID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c *gin.Context, role model.AccessRole) bool {
	return a.IsLowerRoleFn(c, role)
}
