package request

import (
	"github.com/alpacahq/ribbit-backend/apperr"

	"github.com/gin-gonic/gin"
)

// Credentials stores the username and password provided in the request
type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login parses out the username and password in gin's request context, into Credentials
func Login(c *gin.Context) (*Credentials, error) {
	cred := new(Credentials)
	if err := c.ShouldBindJSON(cred); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return cred, nil
}

// ForgotPayload stores the email provided in the request
type ForgotPayload struct {
	Email string `json:"email" binding:"required"`
}

// Forgot parses out the email in gin's request context, into ForgotPayload
func Forgot(c *gin.Context) (*ForgotPayload, error) {
	fgt := new(ForgotPayload)
	if err := c.ShouldBindJSON(fgt); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return fgt, nil
}

// RecoverPasswordPayload stores the data provided in the request
type RecoverPasswordPayload struct {
	Email           string `json:"email" binding:"required"`
	OTP             string `json:"otp" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confrim_password" binding:"required"`
}

// RecoverPassword parses out the data in gin's request context, into RecoverPasswordPayload
func RecoverPassword(c *gin.Context) (*RecoverPasswordPayload, error) {
	rpp := new(RecoverPasswordPayload)
	if err := c.ShouldBindJSON(rpp); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return rpp, nil
}
