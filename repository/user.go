package repository

import (
	"net/http"

	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
)

const notDeleted = "deleted_at is null"

// NewUserRepo returns a new UserRepo instance
func NewUserRepo(db orm.DB, log *zap.Logger) *UserRepo {
	return &UserRepo{db, log}
}

// UserRepo is the client for our user model
type UserRepo struct {
	db  orm.DB
	log *zap.Logger
}

// View returns single user by ID
func (u *UserRepo) View(id int) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, id)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// View returns single user by referral code
func (u *UserRepo) FindByReferralCode(referralCode string) (*model.ReferralCodeVerifyResponse, error) {
	var user = new(model.ReferralCodeVerifyResponse)
	sql := `SELECT "user"."first_name", "user"."last_name", "user"."referral_code", "user"."username"
	FROM "users" AS "user" 
	WHERE ("user"."referral_code" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, referralCode)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByUsername queries for a single user by username
func (u *UserRepo) FindByUsername(username string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, username)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByEmail queries for a single user by email
func (u *UserRepo) FindByEmail(email string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."email" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, email)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByMobile queries for a single user by mobile (and country code)
func (u *UserRepo) FindByMobile(countryCode, mobile string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."country_code" = ? and "user"."mobile" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, countryCode, mobile)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *UserRepo) FindByToken(token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// UpdateLogin updates last login and refresh token for user
func (u *UserRepo) UpdateLogin(user *model.User) error {
	user.UpdateLastLogin() // update user object's last_login field
	_, err := u.db.Model(user).Column("last_login", "token").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}

// List returns list of all users retreivable for the current user, depending on role
func (u *UserRepo) List(qp *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	var users []model.User
	q := u.db.Model(&users).Column("user.*", "Role").Limit(p.Limit).Offset(p.Offset).Where(notDeleted).Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// Update updates user's contact info
func (u *UserRepo) Update(user *model.User) (*model.User, error) {
	_, err := u.db.Model(user).Column(
		"first_name",
		"last_name",
		"username",
		"mobile",
		"country_code",
		"address",
		"account_id",
		"account_number",
		"account_currency",
		"account_status",
		"dob",
		"city",
		"state",
		"country",
		"tax_id_type",
		"tax_id",
		"funding_source",
		"employment_status",
		"investing_experience",
		"public_shareholder",
		"another_brokerage",
		"device_id",
		"profile_completion",
		"bio",
		"facebook_url",
		"twitter_url",
		"instagram_url",
		"public_portfolio",
		"employer_name",
		"occupation",
		"unit_apt",
		"zip_code",
		"stock_symbol",
		"brokerage_firm_name",
		"brokerage_firm_employee_name",
		"brokerage_firm_employee_relationship",
		"shareholder_company_name",
		"avatar",
		"referred_by",
		"watchlist_id",
		"active",
		"verified",
		"updated_at",
	).WherePK().Update()
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
	}
	return user, err
}

// Delete sets deleted_at for a user
func (u *UserRepo) Delete(user *model.User) error {
	user.Delete()
	_, err := u.db.Model(user).Column("deleted_at").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}
