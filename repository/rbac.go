package repository

import (
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/gin-gonic/gin"
)

// NewRBACService creates new RBAC service
func NewRBACService(userRepo model.UserRepo) *RBACService {
	return &RBACService{
		userRepo: userRepo,
	}
}

// RBACService is RBAC application service
type RBACService struct {
	userRepo model.UserRepo
}

// EnforceRole authorizes request by AccessRole
func (s *RBACService) EnforceRole(c *gin.Context, r model.AccessRole) bool {
	return !(c.MustGet("role").(int8) > int8(r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *RBACService) EnforceUser(c *gin.Context, ID int) bool {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	return (c.GetInt("id") == ID) || s.isAdmin(c)
}

func (s *RBACService) isAdmin(c *gin.Context) bool {
	return !(c.MustGet("role").(int8) > int8(model.AdminRole))
}

// AccountCreate performs auth check when creating a new account
// Location admin cannot create accounts, needs to be fixed on EnforceLocation function
func (s *RBACService) AccountCreate(c *gin.Context, roleID int) bool {
	roleCheck := s.EnforceRole(c, model.AccessRole(roleID))
	return roleCheck && s.IsLowerRole(c, model.AccessRole(roleID))
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *RBACService) IsLowerRole(c *gin.Context, r model.AccessRole) bool {
	return !(c.MustGet("role").(int8) >= int8(r))
}
