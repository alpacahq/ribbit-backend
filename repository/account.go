package repository

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"
)

// NewAccountRepo returns an AccountRepo instance
func NewAccountRepo(db orm.DB, log *zap.Logger, secret secret.Service) *AccountRepo {
	return &AccountRepo{db, log, secret}
}

// AccountRepo represents the client for the user table
type AccountRepo struct {
	db     orm.DB
	log    *zap.Logger
	Secret secret.Service
}

// Create creates a new user in our database
func (a *AccountRepo) Create(u *model.User) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "User already exists.")
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return u, nil
}

func encodeToString(max int) string {
	b := make([]byte, max)
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// CreateAndVerify creates a new user in our database, and generates a verification token.
// User active being false until after verification.
func (a *AccountRepo) CreateAndVerify(u *model.User) (*model.Verification, error) {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err == apperr.DB {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "User already exists.")
	}
	u.Password = a.Secret.HashPassword(u.Password)

	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}

	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatal(err)
	}
	u.ReferralCode = strings.ToUpper(re.ReplaceAllString(strings.Split(u.Email, "@")[0]+strconv.Itoa(u.ID), ""))
	if err := a.db.Update(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}

	v := new(model.Verification)
	v.UserID = u.ID
	v.Token = encodeToString(6)
	if err := a.db.Insert(v); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return v, nil
}

func (a *AccountRepo) CreateForgotToken(u *model.User) (*model.Verification, error) {
	v := new(model.Verification)
	v.UserID = u.ID
	v.Token = encodeToString(6)
	if err := a.db.Insert(v); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return v, nil
}

func (a *AccountRepo) CreateNewOTP(u *model.User) (*model.Verification, error) {
	v := new(model.Verification)
	v.UserID = u.ID
	v.Token = encodeToString(6)
	if err := a.db.Insert(v); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return v, nil
}

// CreateWithMobile creates a new user in our database with country code and mobile number
func (a *AccountRepo) CreateWithMobile(u *model.User) error {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err == apperr.DB {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return apperr.DB
	}
	if res.RowsReturned() != 0 && user.Verified == true {
		return apperr.New(http.StatusConflict, "user already exists and is already verified") // user already exists and is already verified
	}
	if res.RowsReturned() != 0 {
		return apperr.New(http.StatusConflict, "user already exists") // user already exists but is not yet verified
	}
	// generate a cryptographically secure random password hash for this user
	u.Password, err = a.Secret.HashRandomPassword()
	if err != nil {
		return apperr.DB
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return apperr.DB
	}
	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatal(err)
	}
	u.ReferralCode = strings.ToUpper(re.ReplaceAllString(strings.Split(u.Email, "@")[0]+strconv.Itoa(u.ID), ""))

	if err := a.db.Update(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return apperr.DB
	}
	return nil
}

// CreateWithMagic creates a new user in our database with Email
func (a *AccountRepo) CreateWithMagic(u *model.User) (int, error) {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE email = ? AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Email)
	if err == apperr.DB {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return 0, apperr.DB
	}
	if res.RowsReturned() != 0 && user.Verified == true {
		return 0, apperr.NewStatus(http.StatusConflict) // user already exists and is already verified
	}
	// if res.RowsReturned() != 0 {
	// 	return 0, apperr.BadRequest // user already exists but is not yet verified
	// }
	// generate a cryptographically secure random password hash for this user
	u.Password, err = a.Secret.HashRandomPassword()
	if err != nil {
		return 0, apperr.DB
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return 0, apperr.DB
	}
	re, err := regexp.Compile(`[^a-z0-9]`)
	if err != nil {
		log.Fatal(err)
	}
	u.ReferralCode = strings.ToUpper(re.ReplaceAllString(strings.Split(u.Email, "@")[0]+strconv.Itoa(u.ID), ""))
	if err := a.db.Update(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return 0, apperr.DB
	}
	return u.ID, nil
}

// ChangePassword changes user's password
func (a *AccountRepo) ChangePassword(u *model.User) error {
	u.Update()
	_, err := a.db.Model(u).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}

// ResetPassword changes user's password
func (a *AccountRepo) ResetPassword(u *model.User) error {
	u.Update()
	u.Password = a.Secret.HashPassword(u.Password)
	_, err := a.db.Model(u).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}

// UpdateAvatar changes user's avatar
func (a *AccountRepo) UpdateAvatar(u *model.User) error {
	u.Update()
	_, err := a.db.Model(u).Column("avatar", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}

// Activate changes user's password
func (a *AccountRepo) Activate(u *model.User) error {
	u.Update()
	_, err := a.db.Model(u).Column("active", "verified", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}

// FindVerificationToken retrieves an existing verification token
func (a *AccountRepo) FindVerificationToken(token string) (*model.Verification, error) {
	var v = new(model.Verification)
	sql := `SELECT * FROM verifications WHERE (token = ? and deleted_at IS NULL)`
	_, err := a.db.QueryOne(v, sql, token)
	if err != nil {
		a.log.Warn("AccountRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "Invalid OTP")
	}
	user := new(model.User)
	user.ID = v.UserID
	user.Active = true
	user.Verified = true
	a.Activate(user)
	return v, nil
}

// FindVerificationTokenByUser retrieves an existing verification token
func (a *AccountRepo) FindVerificationTokenByUser(user *model.User) (*model.Verification, error) {
	var v = new(model.Verification)
	sql := `SELECT * FROM verifications WHERE (user_id = ? and deleted_at IS NULL)`
	_, err := a.db.QueryOne(v, sql, user.ID)
	if err != nil {
		a.log.Warn("AccountRepo Error", zap.String("Error:", err.Error()))
		return nil, nil
		// return nil, apperr.New(http.StatusNotFound, "Invalid OTP")
	}
	return v, nil
}

// DeleteVerificationToken sets deleted_at for an existing verification token
func (a *AccountRepo) DeleteVerificationToken(v *model.Verification) error {
	v.Delete()
	_, err := a.db.Model(v).Column("deleted_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error", zap.Error(err))
		return apperr.DB
	}
	return err
}
