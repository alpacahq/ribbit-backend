package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/repository/auth"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/gin-gonic/gin"
)

// AuthRouter creates new auth http service
func AuthRouter(svc *auth.Service, r *gin.Engine) {
	a := Auth{svc}
	r.POST("/mobile", a.mobile) // mobile: passwordless authentication which handles both the signup scenario and the login scenario
	r.POST("/magic", a.magic)   // magic: magic link authentication which handles both the signup scenario and the login scenario
	r.POST("/signup", a.signup) // email: creates user object
	r.POST("/login", a.login)
	r.POST("/forgot-password", a.forgot)
	r.POST("/recover-password", a.recoverPassword)
	r.GET("/refresh/:token", a.refresh)
	r.GET("/verification/:token", a.verify)                             // email: on verification token submission, mark user as verified and return jwt
	r.POST("/mobile/verify", a.mobileVerify)                            // mobile: on sms code submission, either mark user as verified and return jwt, or update last_login and return jwt
	r.GET("/referral_code/verify/:referral_code", a.referralCodeVerify) // verify referral code
	r.GET("/terms-condition", a.termsCondition)
}

// Auth represents auth http service
type Auth struct {
	svc *auth.Service
}

func (a *Auth) termsCondition(c *gin.Context) {
	c.HTML(http.StatusOK, "terms_conditions.html", gin.H{})
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func (a *Auth) login(c *gin.Context) {
	cred, err := request.Login(c)
	if err != nil {
		return
	}
	fileByte, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Password doesn't match."))
		return
	}
	priv_parsed, err := ParseRsaPrivateKeyFromPemStr(string(fileByte))
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Password doesn't match."))
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(cred.Password)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Password doesn't match."))
		return
	}
	password, err := rsa.DecryptPKCS1v15(rand.Reader, priv_parsed, decoded)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Password doesn't match."))
		return
	}

	r, err := a.svc.Authenticate(c, cred.Email, string(password))
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, r)
}

func (a *Auth) refresh(c *gin.Context) {
	refreshToken := c.Param("token")
	r, err := a.svc.Refresh(c, refreshToken)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, r)
}

func (a *Auth) signup(c *gin.Context) {
	e, err := request.AccountSignup(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}

	fileByte, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}
	priv_parsed, err := ParseRsaPrivateKeyFromPemStr(string(fileByte))
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(e.Password)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}
	password, err := rsa.DecryptPKCS1v15(rand.Reader, priv_parsed, decoded)
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}

	user, err := a.svc.Signup(c, e, string(password))
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)

}

func (a *Auth) verify(c *gin.Context) {
	token := c.Param("token")
	err := a.svc.Verify(c, token)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Auth) forgot(c *gin.Context) {
	body, e := request.Forgot(c)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}
	err := a.svc.Forgot(c, body.Email)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent.",
	})
}

func (a *Auth) recoverPassword(c *gin.Context) {
	body, e := request.RecoverPassword(c)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}

	fileByte, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}
	priv_parsed, err := ParseRsaPrivateKeyFromPemStr(string(fileByte))
	if err != nil {
		apperr.Response(c, apperr.New(http.StatusUnauthorized, "Invalid password format."))
		return
	}

	passwordDecoded, err := base64.StdEncoding.DecodeString(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}
	password, err := rsa.DecryptPKCS1v15(rand.Reader, priv_parsed, passwordDecoded)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}

	confirmPasswordDecoded, err := base64.StdEncoding.DecodeString(body.ConfirmPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}
	ConfirmPassword, err := rsa.DecryptPKCS1v15(rand.Reader, priv_parsed, confirmPasswordDecoded)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters.",
		})
		return
	}

	if string(password) != string(ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password and confirm password doesn't match.",
		})
		return
	}

	errorMessage := a.svc.RecoverPassword(c, body.Email, body.OTP, string(password))
	if errorMessage != nil {
		apperr.Response(c, errorMessage)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent.",
	})
}

// mobile handles a passwordless mobile signup/login
// if user with country_code and mobile already exists, simply return 200
// if user does not exist yet, we attempt to create the new user object, on success 201, otherwise 500
// the client should call /mobile/verify next, if it receives 201 (newly created user object) or 200 (success, and user was previously created)
// we can use this status code in the client to prepare our request object with Signup attribute as true (201) or false (200)
func (a *Auth) mobile(c *gin.Context) {
	m, err := request.Mobile(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	err = a.svc.Mobile(c, m)
	if err != nil {
		if err.Error() == "User already exists." {
			c.JSON(http.StatusOK, gin.H{
				"message": "User already exists.",
			})
			return
		}
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// magic handles a passwordless magic link signup/login
// we can use status code in the client to prepare our request object with Signup attribute as true (201) or false (200)
func (a *Auth) magic(c *gin.Context) {
	m, err := request.Magic(c)
	if err != nil {
		fmt.Println(m)
		apperr.Response(c, err)
		return
	}
	user, err := a.svc.Magic(c, m)
	if err != nil {
		fmt.Println(user)
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// mobileVerify handles the next API call after the previous client call to /mobile
// we mark user verified AND return jwt
func (a *Auth) mobileVerify(c *gin.Context) {
	m, err := request.AccountVerifyMobile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	r, err := a.svc.MobileVerify(c, m.CountryCode, m.Mobile, m.Code, m.Signup)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	c.JSON(http.StatusOK, r)
}

func (a *Auth) referralCodeVerify(c *gin.Context) {
	referralCode := c.Param("referral_code")
	r, err := a.svc.RefVerify(c, referralCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"first_name":    nil,
			"last_name":     nil,
			"referral_code": nil,
		})
		return
	}
	c.JSON(http.StatusOK, r)
}
