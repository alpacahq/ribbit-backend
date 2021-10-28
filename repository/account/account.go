package account

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/platform/structs"
	"github.com/alpacahq/ribbit-backend/request"
	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/gin-gonic/gin"
)

// Service represents the account application service
type Service struct {
	accountRepo model.AccountRepo
	userRepo    model.UserRepo
	rbac        model.RBACService
	secret      secret.Service
}

// NewAccountService creates a new account application service
func NewAccountService(userRepo model.UserRepo, accountRepo model.AccountRepo, rbac model.RBACService, secret secret.Service) *Service {
	return &Service{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		rbac:        rbac,
		secret:      secret,
	}
}

// Create creates a new user account
func (s *Service) Create(c *gin.Context, u *model.User) error {
	if !s.rbac.AccountCreate(c, u.RoleID) {
		return apperr.New(http.StatusForbidden, "Forbidden")
	}
	u.Password = s.secret.HashPassword(u.Password)
	u, err := s.accountRepo.Create(u)
	return err
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(c *gin.Context, oldPass, newPass string, id int) error {
	if !s.rbac.EnforceUser(c, id) {
		return apperr.New(http.StatusForbidden, "Forbidden")
	}
	u, err := s.userRepo.View(id)
	if err != nil {
		return err
	}
	if !s.secret.HashMatchesPassword(u.Password, oldPass) {
		return apperr.New(http.StatusBadGateway, "old password is not correct")
	}
	u.Password = s.secret.HashPassword(newPass)
	return s.accountRepo.ChangePassword(u)
}

// UpdateAvatar changes user's avatar
func (s *Service) UpdateAvatar(c *gin.Context, newAvatar string, id int) error {
	if !s.rbac.EnforceUser(c, id) {
		return apperr.New(http.StatusForbidden, "Forbidden")
	}
	u, err := s.userRepo.View(id)
	if err != nil {
		return err
	}
	u.Avatar = newAvatar
	return s.accountRepo.UpdateAvatar(u)
}

// GetProfile gets user's profile
func (s *Service) GetProfile(c *gin.Context, id int) *model.User {
	if !s.rbac.EnforceUser(c, id) {
		return nil
	}
	u, err := s.userRepo.View(id)
	if err != nil {
		return nil
	}

	return u
}

// UpdateProfile updated user's profile
func (s *Service) UpdateProfile(c *gin.Context, update *request.Update) (*model.User, error) {
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
