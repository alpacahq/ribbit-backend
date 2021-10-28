package repository

import (
	"fmt"
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"
)

// NewAssetRepo returns an AssetRepo instance
func NewAssetRepo(db orm.DB, log *zap.Logger, secret secret.Service) *AssetRepo {
	return &AssetRepo{db, log, secret}
}

// AssetRepo represents the client for the user table
type AssetRepo struct {
	db     orm.DB
	log    *zap.Logger
	Secret secret.Service
}

// Create creates a new asset in our database.
func (a *AssetRepo) CreateOrUpdate(ass *model.Asset) (*model.Asset, error) {
	_asset := new(model.Asset)
	sql := `SELECT id FROM assets WHERE symbol = ?`
	res, err := a.db.Query(_asset, sql, ass.Symbol)
	if err == apperr.DB {
		a.log.Error("AssetRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		// update..
		fmt.Println("updating...")
		_, err := a.db.Model(ass).Column(
			"class",
			"exchange",
			"name",
			"status",
			"tradable",
			"marginable",
			"shortable",
			"easy_to_borrow",
			"fractionable",
			"is_watchlisted",
			"updated_at",
		).WherePK().Update()
		if err != nil {
			a.log.Warn("AssetRepo Error: ", zap.Error(err))
			return nil, err
		}
		return ass, nil
	} else {
		// create
		fmt.Println("creating...")
		if err := a.db.Insert(ass); err != nil {
			a.log.Warn("AssetRepo error: ", zap.Error(err))
			return nil, apperr.DB
		}
	}
	return ass, nil
}

// UpdateAsset changes user's avatar
func (a *AssetRepo) UpdateAsset(u *model.Asset) error {
	_, err := a.db.Model(u).Column(
		"class",
		"exchange",
		"name",
		"status",
		"tradable",
		"marginable",
		"shortable",
		"easy_to_borrow",
		"fractionable",
		"is_watchlisted",
		"updated_at",
	).WherePK().Update()
	if err != nil {
		a.log.Warn("AssetRepo Error: ", zap.Error(err))
	}
	return err
}

// SearchAssets changes user's avatar
func (a *AssetRepo) Search(query string) ([]model.Asset, error) {
	var exactAsset model.Asset
	var assets []model.Asset
	sql := `SELECT * FROM assets WHERE LOWER(symbol) = LOWER(?) LIMIT 1`
	_, err := a.db.QueryOne(&exactAsset, sql, query, query, query)
	if err != nil {
		a.log.Warn("AssetRepo Error", zap.String("Error:", err.Error()))
	}

	sql2 := `SELECT * FROM assets WHERE symbol ILIKE ? || '%' OR name ILIKE ? || '%' ORDER BY symbol ASC LIMIT 50`
	_, err2 := a.db.Query(&assets, sql2, query, query, query)
	if err2 != nil {
		a.log.Warn("AssetRepo Error", zap.String("Error:", err2.Error()))
		return assets, apperr.New(http.StatusNotFound, "404 not found")
	}

	if err == nil {
		fmt.Println(exactAsset)
		assets = append([]model.Asset{exactAsset}, findAndDelete(assets, exactAsset)...)
	}

	return assets, nil
}

func findAndDelete(s []model.Asset, item model.Asset) []model.Asset {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}
