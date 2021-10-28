package request

import (
	"net/http"

	apperr "github.com/alpacahq/ribbit-backend/apperr"
	model "github.com/alpacahq/ribbit-backend/model"

	"github.com/gin-gonic/gin"
)

// RegisterAdmin contains admin registration request
type RegisterAdmin struct {
	FirstName                         string `json:"first_name" binding:"required"`
	LastName                          string `json:"last_name" binding:"required"`
	Username                          string `json:"username" binding:"required,min=3,alphanum"`
	Password                          string `json:"password" binding:"required,min=8"`
	Email                             string `json:"email" binding:"required,email"`
	RoleID                            int    `json:"role_id" binding:"required"`
	AccountID                         string `json:"account_id"`
	AccountNumber                     string `json:"account_number"`
	AccountCurrency                   string `json:"account_currency"`
	AccountStatus                     string `json:"account_status"`
	DOB                               string `json:"dob"`
	City                              string `json:"city"`
	State                             string `json:"state"`
	Country                           string `json:"country"`
	TaxIDType                         string `json:"tax_id_type"`
	TaxID                             string `json:"tax_id"`
	FundingSource                     string `json:"funding_source"`
	EmploymentStatus                  string `json:"employment_status"`
	InvestingExperience               string `json:"investing_experience"`
	PublicShareholder                 string `json:"public_shareholder"`
	AnotherBrokerage                  string `json:"another_brokerage"`
	DeviceID                          string `json:"device_id"`
	ProfileCompletion                 string `json:"profile_completion"`
	BIO                               string `json:"bio"`
	FacebookURL                       string `json:"facebook_url"`
	TwitterURL                        string `json:"twitter_url"`
	InstagramURL                      string `json:"instagram_url"`
	PublicPortfolio                   string `json:"public_portfolio"`
	EmployerName                      string `json:"employer_name"`
	Occupation                        string `json:"occupation"`
	UnitApt                           string `json:"unit_apt"`
	ZipCode                           string `json:"zip_code"`
	StockSymbol                       string `json:"stock_symbol"`
	BrokerageFirmName                 string `json:"brokerage_firm_name"`
	BrokerageFirmEmployeeName         string `json:"brokerage_firm_employee_name"`
	BrokerageFirmEmployeeRelationship string `json:"brokerage_firm_employee_relationship"`
	ShareholderCompanyName            string `json:"shareholder_company_name"`
	Avatar                            string `json:"avatar"`
	ReferralCode                      string `json:"referral_code"`
	ReferredBy                        string `json:"referred_by"`
}

// AccountCreate validates account creation request
func AccountCreate(c *gin.Context) (*RegisterAdmin, error) {
	var r RegisterAdmin
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	if r.RoleID < int(model.SuperAdminRole) || r.RoleID > int(model.UserRole) {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, apperr.New(http.StatusBadRequest, "Couldn't create account.")
	}
	return &r, nil
}

// Password contains password change request
type Password struct {
	ID          int    `json:"-"`
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// PasswordChange validates password change request
func PasswordChange(c *gin.Context) (*Password, error) {
	var p Password
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	p.ID = id
	return &p, nil
}

// Update contains password change request
type Update struct {
	ID                                int     `json:"id"`
	FirstName                         *string `json:"first_name"`
	LastName                          *string `json:"last_name"`
	Username                          *string `json:"username"`
	Password                          *string `json:"-"`
	Email                             *string `json:"email"`
	Mobile                            *string `json:"mobile"`
	CountryCode                       *string `json:"country_code"`
	Address                           *string `json:"address"`
	AccountID                         *string `json:"account_id"`
	AccountNumber                     *string `json:"account_number"`
	AccountCurrency                   *string `json:"account_currency"`
	AccountStatus                     *string `json:"account_status"`
	DOB                               *string `json:"dob"`
	City                              *string `json:"city"`
	State                             *string `json:"state"`
	Country                           *string `json:"country"`
	TaxIDType                         *string `json:"tax_id_type"`
	TaxID                             *string `json:"tax_id"`
	FundingSource                     *string `json:"funding_source"`
	EmploymentStatus                  *string `json:"employment_status"`
	InvestingExperience               *string `json:"investing_experience"`
	PublicShareholder                 *string `json:"public_shareholder"`
	AnotherBrokerage                  *string `json:"another_brokerage"`
	DeviceID                          *string `json:"device_id"`
	ProfileCompletion                 *string `json:"profile_completion"`
	BIO                               *string `json:"bio"`
	FacebookURL                       *string `json:"facebook_url"`
	TwitterURL                        *string `json:"twitter_url"`
	InstagramURL                      *string `json:"instagram_url"`
	PublicPortfolio                   *string `json:"public_portfolio"`
	EmployerName                      *string `json:"employer_name"`
	Occupation                        *string `json:"occupation"`
	UnitApt                           *string `json:"unit_apt"`
	ZipCode                           *string `json:"zip_code"`
	StockSymbol                       *string `json:"stock_symbol"`
	BrokerageFirmName                 *string `json:"brokerage_firm_name"`
	BrokerageFirmEmployeeName         *string `json:"brokerage_firm_employee_name"`
	BrokerageFirmEmployeeRelationship *string `json:"brokerage_firm_employee_relationship"`
	ShareholderCompanyName            *string `json:"shareholder_company_name"`
	Avatar                            *string `json:"avatar"`
	ReferredBy                        *string `json:"referred_by"`
	WatchlistID                       *string `json:"watchlist_id"`
}

// UpdateProfile updates user's profile
func UpdateProfile(c *gin.Context) (*Update, error) {
	var p Update
	id, _ := c.Get("id")
	if err := c.ShouldBindJSON(&p); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	p.ID = id.(int)
	return &p, nil
}
