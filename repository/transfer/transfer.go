package transfer

import (
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"
)

// NewAuthService creates new auth service
func NewTransferService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT, db orm.DB, log *zap.Logger) *Service {
	return &Service{userRepo, accountRepo, jwt, db, log}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
	db          orm.DB
	log         *zap.Logger
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}
