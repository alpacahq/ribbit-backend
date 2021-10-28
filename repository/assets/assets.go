package assets

import (
	"fmt"

	"github.com/alpacahq/ribbit-backend/model"
	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(fmt.Errorf("unexpected error while initializing plaid client %w", err))
	}
}

// NewAuthService creates new auth service
func NewAssetsService(userRepo model.UserRepo, accountRepo model.AccountRepo, assetRepo model.AssetsRepo, jwt JWT, db orm.DB, log *zap.Logger) *Service {
	return &Service{userRepo, assetRepo, accountRepo, jwt, db, log}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	assetRepo   model.AssetsRepo
	accountRepo model.AccountRepo
	jwt         JWT
	db          orm.DB
	log         *zap.Logger
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

// SearchAssets changes user's avatar
func (a *Service) SearchAssets(query string) ([]model.Asset, error) {
	assets, err := a.assetRepo.Search(query)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
