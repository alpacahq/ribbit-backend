package request

import (
	"github.com/alpacahq/ribbit-backend/apperr"

	"github.com/gin-gonic/gin"
)

// EmailSignup contains the user signup request
type EmailSignup struct {
	Email    string `json:"email" binding:"required,min=3,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// AccountSignup validates user signup request
func AccountSignup(c *gin.Context) (*EmailSignup, error) {
	var r EmailSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return &r, nil
}

// MobileSignup contains the user signup request with a mobile number
type MobileSignup struct {
	CountryCode string `json:"country_code" binding:"required,min=2"`
	Mobile      string `json:"mobile" binding:"required"`
}

// Mobile validates user signup request via mobile
func Mobile(c *gin.Context) (*MobileSignup, error) {
	var r MobileSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return &r, nil
}

// MagicSignup contains the user signup request with a mobile number
type MagicSignup struct {
	Email string `json:"email" binding:"required,min=3,email"`
}

// Magic validates user signup request via mobile
func Magic(c *gin.Context) (*MagicSignup, error) {
	var r MagicSignup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return &r, nil
}

// MobileVerify contains the user's mobile verification country code, mobile number and verification code
type MobileVerify struct {
	CountryCode string `json:"country_code" binding:"required,min=2"`
	Mobile      string `json:"mobile" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Signup      bool   `json:"signup" binding:"required"`
}

// AccountVerifyMobile validates user mobile verification
func AccountVerifyMobile(c *gin.Context) (*MobileVerify, error) {
	var r MobileVerify
	if err := c.ShouldBindJSON(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

type ReferralVerify struct {
	ReferralCode string `json:"referral_code" binding:"required"`
}

// ReferralCodeVerify verifies referral code
func ReferralCodeVerify(c *gin.Context) (*ReferralVerify, error) {
	var r ReferralVerify
	if err := c.ShouldBindJSON(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
