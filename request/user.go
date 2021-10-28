package request

import (
	"github.com/alpacahq/ribbit-backend/apperr"

	"github.com/gin-gonic/gin"
)

// UpdateUser contains user update data from json request
type UpdateUser struct {
	ID                                int     `json:"-"`
	FirstName                         *string `json:"first_name,omitempty" binding:"omitempty,min=2"`
	LastName                          *string `json:"last_name,omitempty" binding:"omitempty,min=2"`
	Mobile                            *string `json:"mobile,omitempty"`
	Phone                             *string `json:"phone,omitempty"`
	Address                           *string `json:"address,omitempty"`
	AccountID                         *string `json:"account_id,omitempty"`
	AccountNumber                     *string `json:"account_number,omitempty"`
	AccountCurrency                   *string `json:"account_currency,omitempty"`
	AccountStatus                     *string `json:"account_status,omitempty"`
	DOB                               *string `json:"dob,omitempty"`
	City                              *string `json:"city,omitempty"`
	State                             *string `json:"state,omitempty"`
	Country                           *string `json:"country,omitempty"`
	TaxIDType                         *string `json:"tax_id_type,omitempty"`
	TaxID                             *string `json:"tax_id,omitempty"`
	FundingSource                     *string `json:"funding_source,omitempty"`
	EmploymentStatus                  *string `json:"employment_status"`
	InvestingExperience               *string `json:"investing_experience,omitempty"`
	PublicShareholder                 *string `json:"public_shareholder,omitempty"`
	AnotherBrokerage                  *string `json:"another_brokerage,omitempty"`
	DeviceID                          *string `json:"device_id,omitempty"`
	ProfileCompletion                 *string `json:"profile_completion,omitempty"`
	BIO                               *string `json:"bio,omitempty"`
	FacebookURL                       *string `json:"facebook_url,omitempty"`
	TwitterURL                        *string `json:"twitter_url,omitempty"`
	InstagramURL                      *string `json:"instagram_url,omitempty"`
	PublicPortfolio                   *string `json:"public_portfolio,omitempty"`
	EmployerName                      *string `json:"employer_name,omitempty"`
	Occupation                        *string `json:"occupation,omitempty"`
	UnitApt                           *string `json:"unit_apt,omitempty"`
	ZipCode                           *string `json:"zip_code,omitempty"`
	StockSymbol                       *string `json:"stock_symbol,omitempty"`
	BrokerageFirmName                 *string `json:"brokerage_firm_name,omitempty"`
	BrokerageFirmEmployeeName         *string `json:"brokerage_firm_employee_name,omitempty"`
	BrokerageFirmEmployeeRelationship *string `json:"brokerage_firm_employee_relationship,omitempty"`
	ShareholderCompanyName            *string `json:"shareholder_company_name,omitempty"`
	Avatar                            *string `json:"avatar,omitempty"`
	ReferredBy                        *string `json:"referred_by,omitempty"`
	ReferralCode                      *string `json:"referral_code,omitempty"`
}

// UserUpdate validates user update request
func UserUpdate(c *gin.Context) (*UpdateUser, error) {
	var u UpdateUser
	id, err := ID(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(&u); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	u.ID = id
	return &u, nil
}
