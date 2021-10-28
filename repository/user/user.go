package user

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/platform/query"
	"github.com/alpacahq/ribbit-backend/repository/platform/structs"

	"github.com/gin-gonic/gin"
)

// NewUserService create a new user application service
func NewUserService(userRepo model.UserRepo, auth model.AuthService, rbac model.RBACService) *Service {
	return &Service{
		userRepo: userRepo,
		auth:     auth,
		rbac:     rbac,
	}
}

// Service represents the user application service
type Service struct {
	userRepo model.UserRepo
	auth     model.AuthService
	rbac     model.RBACService
}

// List returns list of users
func (s *Service) List(c *gin.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.userRepo.List(q, p)
}

// View returns single user
func (s *Service) View(c *gin.Context, id int) (*model.User, error) {
	if !s.rbac.EnforceUser(c, id) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	return s.userRepo.View(id)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName *string
	LastName  *string
	Mobile    *string
	Phone     *string
	Address   *string
}

// Update updates user's contact information
func (s *Service) Update(c *gin.Context, update *Update) (*model.User, error) {
	if !s.rbac.EnforceUser(c, update.ID) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	u, err := s.userRepo.View(update.ID)
	if err != nil {
		return nil, err
	}
	structs.Merge(u, update)
	return s.userRepo.Update(u)
}

// Delete deletes a user
func (s *Service) Delete(c *gin.Context, id int) error {
	u, err := s.userRepo.View(id)
	if err != nil {
		return err
	}
	if !s.rbac.IsLowerRole(c, u.Role.AccessLevel) {
		return apperr.New(http.StatusForbidden, "Forbidden")
	}
	u.Delete()
	return s.userRepo.Delete(u)
}
