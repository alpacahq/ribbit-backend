package mockdb

import (
	"github.com/alpacahq/ribbit-backend/model"
)

// Account database mock
type Account struct {
	ActivateFn                    func(*model.User) error
	CreateFn                      func(*model.User) (*model.User, error)
	CreateAndVerifyFn             func(*model.User) (*model.Verification, error)
	CreateWithMobileFn            func(*model.User) error
	CreateForgotTokenFn           func(*model.User) (*model.Verification, error)
	CreateNewOTPFn                func(*model.User) (*model.Verification, error)
	CreateWithMagicFn             func(*model.User) (int, error)
	ChangePasswordFn              func(*model.User) error
	ResetPasswordFn               func(*model.User) error
	UpdateAvatarFn                func(*model.User) error
	FindVerificationTokenFn       func(string) (*model.Verification, error)
	FindVerificationTokenByUserFn func(*model.User) (*model.Verification, error)
	DeleteVerificationTokenFn     func(*model.Verification) error
}

func (a *Account) Activate(usr *model.User) error {
	return a.ActivateFn(usr)
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

func (a *Account) CreateForgotToken(usr *model.User) (*model.Verification, error) {
	return a.CreateForgotTokenFn(usr)
}

func (a *Account) CreateNewOTP(usr *model.User) (*model.Verification, error) {
	return a.CreateNewOTPFn(usr)
}

func (a *Account) CreateWithMagic(usr *model.User) (int, error) {
	return a.CreateWithMagicFn(usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(usr *model.User) error {
	return a.ChangePasswordFn(usr)
}

func (a *Account) UpdateAvatar(usr *model.User) error {
	return a.UpdateAvatarFn(usr)
}

func (a *Account) ResetPassword(usr *model.User) error {
	return a.ResetPasswordFn(usr)
}

// FindVerificationToken mock
func (a *Account) FindVerificationToken(token string) (*model.Verification, error) {
	return a.FindVerificationTokenFn(token)
}

func (a *Account) FindVerificationTokenByUser(usr *model.User) (*model.Verification, error) {
	return a.FindVerificationTokenByUserFn(usr)
}

// DeleteVerificationToken mock
func (a *Account) DeleteVerificationToken(v *model.Verification) error {
	return a.DeleteVerificationTokenFn(v)
}
