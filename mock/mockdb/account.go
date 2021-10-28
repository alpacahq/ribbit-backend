package mockdb

import (
	"github.com/alpacahq/ribbit-backend/model"
)

// Account database mock
type Account struct {
	CreateFn                  func(*model.User) (*model.User, error)
	CreateAndVerifyFn         func(*model.User) (*model.Verification, error)
	CreateWithMobileFn        func(*model.User) error
	ChangePasswordFn          func(*model.User) error
	FindVerificationTokenFn   func(string) (*model.Verification, error)
	DeleteVerificationTokenFn func(*model.Verification) error
}

// Create mock
func (a *Account) Create(usr *model.User) (*model.User, error) {
	return a.CreateFn(usr)
}

// CreateAndVerify mock
func (a *Account) CreateAndVerify(usr *model.User) (*model.Verification, error) {
	return a.CreateAndVerifyFn(usr)
}

// CreateWithMobile mock
func (a *Account) CreateWithMobile(usr *model.User) error {
	return a.CreateWithMobileFn(usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(usr *model.User) error {
	return a.ChangePasswordFn(usr)
}

// FindVerificationToken mock
func (a *Account) FindVerificationToken(token string) (*model.Verification, error) {
	return a.FindVerificationTokenFn(token)
}

// DeleteVerificationToken mock
func (a *Account) DeleteVerificationToken(v *model.Verification) error {
	return a.DeleteVerificationTokenFn(v)
}
