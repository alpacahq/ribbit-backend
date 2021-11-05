package mockdb

import (
	"github.com/alpacahq/ribbit-backend/model"
)

// User database mock
type User struct {
	ViewFn               func(int) (*model.User, error)
	FindByReferralCodeFn func(string) (*model.ReferralCodeVerifyResponse, error)
	FindByUsernameFn     func(string) (*model.User, error)
	FindByEmailFn        func(string) (*model.User, error)
	FindByMobileFn       func(string, string) (*model.User, error)
	FindByTokenFn        func(string) (*model.User, error)
	UpdateLoginFn        func(*model.User) error
	ListFn               func(*model.ListQuery, *model.Pagination) ([]model.User, error)
	DeleteFn             func(*model.User) error
	UpdateFn             func(*model.User) (*model.User, error)
}

// View mock
func (u *User) View(id int) (*model.User, error) {
	return u.ViewFn(id)
}

// FindByReferralCode mock
func (u *User) FindByReferralCode(username string) (*model.ReferralCodeVerifyResponse, error) {
	return u.FindByReferralCodeFn(username)
}

// FindByUsername mock
func (u *User) FindByUsername(username string) (*model.User, error) {
	return u.FindByUsernameFn(username)
}

// FindByEmail mock
func (u *User) FindByEmail(email string) (*model.User, error) {
	return u.FindByEmailFn(email)
}

// FindByMobile mock
func (u *User) FindByMobile(countryCode, mobile string) (*model.User, error) {
	return u.FindByMobileFn(countryCode, mobile)
}

// FindByToken mock
func (u *User) FindByToken(token string) (*model.User, error) {
	return u.FindByTokenFn(token)
}

// UpdateLogin mock
func (u *User) UpdateLogin(usr *model.User) error {
	return u.UpdateLoginFn(usr)
}

// List mock
func (u *User) List(lq *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	return u.ListFn(lq, p)
}

// Delete mock
func (u *User) Delete(usr *model.User) error {
	return u.DeleteFn(usr)
}

// Update mock
func (u *User) Update(usr *model.User) (*model.User, error) {
	return u.UpdateFn(usr)
}
