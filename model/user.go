package model

import (
	"time"
)

func init() {
	Register(&User{})
}

// User represents user domain model
type User struct {
	Base
	ID                                int        `json:"id"`
	FirstName                         string     `json:"first_name"`
	LastName                          string     `json:"last_name"`
	Username                          string     `json:"username"`
	Password                          string     `json:"-"`
	Email                             string     `json:"email"`
	Mobile                            string     `json:"mobile"`
	CountryCode                       string     `json:"country_code"`
	Address                           string     `json:"address"`
	LastLogin                         *time.Time `json:"last_login,omitempty"`
	Verified                          bool       `json:"verified"`
	Active                            bool       `json:"active"`
	Token                             string     `json:"-"`
	Role                              *Role      `json:"role,omitempty"`
	RoleID                            int        `json:"-"`
	AccountID                         string     `json:"account_id"`
	AccountNumber                     string     `json:"account_number"`
	AccountCurrency                   string     `json:"account_currency"`
	AccountStatus                     string     `json:"account_status"`
	DOB                               string     `json:"dob"`
	City                              string     `json:"city"`
	State                             string     `json:"state"`
	Country                           string     `json:"country"`
	TaxIDType                         string     `json:"tax_id_type"`
	TaxID                             string     `json:"tax_id"`
	FundingSource                     string     `json:"funding_source"`
	EmploymentStatus                  string     `json:"employment_status"`
	InvestingExperience               string     `json:"investing_experience"`
	PublicShareholder                 string     `json:"public_shareholder"`
	AnotherBrokerage                  string     `json:"another_brokerage"`
	DeviceID                          string     `json:"device_id"`
	ProfileCompletion                 string     `json:"profile_completion"`
	BIO                               string     `json:"bio"`
	FacebookURL                       string     `json:"facebook_url"`
	TwitterURL                        string     `json:"twitter_url"`
	InstagramURL                      string     `json:"instagram_url"`
	PublicPortfolio                   string     `json:"public_portfolio"`
	EmployerName                      string     `json:"employer_name"`
	Occupation                        string     `json:"occupation"`
	UnitApt                           string     `json:"unit_apt"`
	ZipCode                           string     `json:"zip_code"`
	StockSymbol                       string     `json:"stock_symbol"`
	BrokerageFirmName                 string     `json:"brokerage_firm_name"`
	BrokerageFirmEmployeeName         string     `json:"brokerage_firm_employee_name"`
	BrokerageFirmEmployeeRelationship string     `json:"brokerage_firm_employee_relationship"`
	ShareholderCompanyName            string     `json:"shareholder_company_name"`
	Avatar                            string     `json:"avatar"`
	ReferredBy                        string     `json:"referred_by"`
	ReferralCode                      string     `json:"referral_code"`
	WatchlistID                       string     `json:"watchlist_id"`
	PerAccountLimit                   float64    `json:"per_account_limit"`
}

// ReferralCodeVerifyResponse
type ReferralCodeVerifyResponse struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	ReferralCode string `json:"referral_code"`
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin() {
	t := time.Now()
	u.LastLogin = &t
}

// Delete updates the deleted_at field
func (u *User) Delete() {
	t := time.Now()
	u.DeletedAt = &t
}

// Update updates the updated_at field
func (u *User) Update() {
	t := time.Now()
	u.UpdatedAt = t
}

// UserRepo represents user database interface (the repository)
type UserRepo interface {
	View(int) (*User, error)
	FindByUsername(string) (*User, error)
	FindByReferralCode(string) (*ReferralCodeVerifyResponse, error)
	FindByEmail(string) (*User, error)
	FindByMobile(string, string) (*User, error)
	FindByToken(string) (*User, error)
	UpdateLogin(*User) error
	List(*ListQuery, *Pagination) ([]User, error)
	Update(*User) (*User, error)
	Delete(*User) error
}

// AccountRepo represents account database interface (the repository)
type AccountRepo interface {
	Create(*User) (*User, error)
	CreateAndVerify(*User) (*Verification, error)
	CreateForgotToken(*User) (*Verification, error)
	CreateNewOTP(*User) (*Verification, error)
	CreateWithMobile(*User) error
	CreateWithMagic(*User) (int, error)
	ResetPassword(*User) error
	ChangePassword(*User) error
	UpdateAvatar(*User) error
	Activate(*User) error
	FindVerificationToken(string) (*Verification, error)
	FindVerificationTokenByUser(*User) (*Verification, error)
	DeleteVerificationToken(*Verification) error
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID       int
	Username string
	Email    string
	Role     AccessRole
}
