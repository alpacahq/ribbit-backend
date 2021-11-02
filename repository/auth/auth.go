package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	mag "github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/client"
	"github.com/magiclabs/magic-admin-go/token"
	"github.com/rs/xid"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/magic"
	"github.com/alpacahq/ribbit-backend/mail"
	"github.com/alpacahq/ribbit-backend/mobile"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/request"
	"github.com/alpacahq/ribbit-backend/secret"

	shortuuid "github.com/lithammer/shortuuid/v3"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT, m mail.Service, mob mobile.Service, mag magic.Service) *Service {
	return &Service{userRepo, accountRepo, jwt, m, mob, mag}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
	m           mail.Service
	mob         mobile.Service
	mag         magic.Service
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c context.Context, email, password string) (*model.LoginResponseWithToken, error) {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, apperr.New(http.StatusUnauthorized, "Invalid credentials. Please check and submit again.")
	}
	if !secret.New().HashMatchesPassword(u.Password, password) {
		return nil, apperr.New(http.StatusUnauthorized, "Invalid credentials. Please check and submit again.")
	}
	// user must be active and verified. Active is enabled/disabled by superadmin user. Verified depends on user verifying via /verification/:token or /mobile/verify
	// if !u.Active || !u.Verified {
	// 	return nil, apperr.New(http.StatusUnauthorized, "User already exists.")
	// }
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.New(http.StatusUnauthorized, "Invalid credentials. Please check and submit again.")
	}
	u.UpdateLastLogin()
	u.Token = xid.New().String()
	if err := s.userRepo.UpdateLogin(u); err != nil {
		return nil, err
	}

	response := &model.LoginResponseWithToken{
		Token:        token,
		Expires:      expire,
		RefreshToken: u.Token,
		User:         *u,
	}

	if !u.Active {
		v, _ := s.accountRepo.FindVerificationTokenByUser(u)
		if v != nil {
			err = s.accountRepo.DeleteVerificationToken(v)
			if err != nil {
				return response, err
			}
			v2, err := s.accountRepo.CreateNewOTP(u)
			if err != nil { // user exists
				return response, apperr.New(http.StatusInternalServerError, "Invalid credentials. Please check and submit again.")
			}

			err = s.m.SendVerificationEmail(email, v2)
			if err != nil {
				return response, err
			}
		}
	}

	return response, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (s *Service) Refresh(c context.Context, refreshToken string) (*model.RefreshToken, error) {
	user, err := s.userRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// this is our re-generated JWT
	token, expire, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, apperr.Generic
	}
	return &model.RefreshToken{
		Token:   token,
		Expires: expire,
	}, nil
}

// Verify verifies the (verification) token and deletes it
func (s *Service) Verify(c context.Context, token string) error {
	v, err := s.accountRepo.FindVerificationToken(token)
	if err != nil {
		return err
	}
	err = s.accountRepo.DeleteVerificationToken(v)
	if err != nil {
		return err
	}
	return nil
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

// Send pin for forgot password
func (s *Service) Forgot(c *gin.Context, email string) error {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil { // user exists
		return apperr.New(http.StatusNotFound, "User doesn't exist.")
	}
	v, err := s.accountRepo.CreateForgotToken(u)
	if err != nil { // user exists
		return apperr.New(http.StatusInternalServerError, "Failed to generate verification process.")
	}

	err = s.m.SendForgotVerificationEmail(email, v)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return apperr.New(http.StatusOK, "OTP sent.")
}

// Verify OTP and recover password
func (s *Service) RecoverPassword(c *gin.Context, email string, otp string, password string) error {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil { // user exists
		return apperr.New(http.StatusNotFound, "User doesn't exist.")
	}

	v, err := s.accountRepo.FindVerificationToken(otp)
	if err != nil {
		return err
	}
	err = s.accountRepo.DeleteVerificationToken(v)
	if err != nil {
		return err
	}

	u.Password = password
	if err := s.accountRepo.ResetPassword(u); err != nil {
		return apperr.New(http.StatusInternalServerError, "Failed to change password, please try again.")
	}

	return apperr.New(http.StatusAccepted, "Password changed.")
}

// RefVerify verifies the referral code
func (s *Service) RefVerify(c context.Context, referralCode string) (*model.ReferralCodeVerifyResponse, error) {
	u, err := s.userRepo.FindByReferralCode(referralCode)
	if err != nil {
		return nil, err
	} else {
		return u, nil
	}
}

// MobileVerify verifies the mobile verification code, i.e. (6-digit) code
func (s *Service) MobileVerify(c context.Context, countryCode, mobile, code string, signup bool) (*model.AuthToken, error) {
	// send code to twilio
	err := s.mob.CheckCode(countryCode, mobile, code)
	if err != nil {
		return nil, err
	}
	u, err := s.userRepo.FindByMobile(countryCode, mobile)
	if err != nil {
		return nil, err
	}
	if signup { // signup case, make user verified and active
		u.Verified = true
		u.Active = true
	} else { // login case, update user's last_login attribute
		u.UpdateLastLogin()
	}
	u, err = s.userRepo.Update(u)
	if err != nil {
		return nil, err
	}

	// generate jwt and return
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.New(http.StatusUnauthorized, "Unauthorized")
	}
	u.UpdateLastLogin()
	u.Token = xid.New().String()
	if err := s.userRepo.UpdateLogin(u); err != nil {
		return nil, err
	}
	return &model.AuthToken{
		Token:        token,
		Expires:      expire,
		RefreshToken: u.Token,
	}, nil
}

// User returns user data stored in jwt token
func (s *Service) User(c *gin.Context) *model.AuthUser {
	id := c.GetInt("id")
	user := c.GetString("username")
	email := c.GetString("email")
	role := c.MustGet("role").(int8)
	return &model.AuthUser{
		ID:       id,
		Username: user,
		Email:    email,
		Role:     model.AccessRole(role),
	}
}

// Signup returns any error from creating a new user in our database
func (s *Service) Signup(c *gin.Context, e *request.EmailSignup, password string) (*model.LoginResponseWithToken, error) {
	_, err := s.userRepo.FindByEmail(e.Email)
	if err == nil { // user already exists
		return nil, apperr.New(http.StatusConflict, "User already exists.")
	}
	u := shortuuid.New()
	v, err := s.accountRepo.CreateAndVerify(&model.User{Email: e.Email, Password: password, ReferralCode: u})
	if err != nil {
		return nil, err
	}
	err = s.m.SendVerificationEmail(e.Email, v)
	if err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	newUser, err2 := s.userRepo.View(v.UserID)
	if err2 == nil { // user already exists
		fmt.Println(newUser)
		// user must be active and verified. Active is enabled/disabled by superadmin user. Verified depends on user verifying via /verification/:token or /mobile/verify
		// if !newUser.Active || !newUser.Verified {
		// 	return nil, apperr.Unauthorized
		// }
		token, expire, err := s.jwt.GenerateToken(newUser)
		if err != nil {
			return nil, apperr.New(http.StatusUnauthorized, "Unauthorized")
		}
		newUser.UpdateLastLogin()
		newUser.Token = xid.New().String()
		if err := s.userRepo.UpdateLogin(newUser); err != nil {
			return nil, err
		}
		return &model.LoginResponseWithToken{
			Token:        token,
			Expires:      expire,
			RefreshToken: newUser.Token,
			User:         *newUser,
		}, nil
	}

	return nil, apperr.NewStatus(http.StatusBadRequest)
}

// Mobile returns any error from creating a new user in our database with a mobile number
func (s *Service) Mobile(c *gin.Context, m *request.MobileSignup) error {
	// find by countryCode and mobile
	_, err := s.userRepo.FindByMobile(m.CountryCode, m.Mobile)
	if err == nil { // user already exists
		return apperr.New(http.StatusConflict, "User already exists.")
	}
	// create and verify
	user := &model.User{
		CountryCode: m.CountryCode,
		Mobile:      m.Mobile,
	}
	err = s.accountRepo.CreateWithMobile(user)
	if err != nil {
		return err
	}
	// generate sms token
	err = s.mob.GenerateSMSToken(m.CountryCode, m.Mobile)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return nil
}

// Magic returns any error from creating a new user in our database with a magic link
func (s *Service) Magic(c *gin.Context, m *request.MagicSignup) (*model.LoginResponseWithToken, error) {
	// Validate magic token

	tkn := c.Request.Header.Get("Authorization")

	authBearer := "Bearer"
	if tkn == "" {
		return nil, apperr.New(http.StatusUnauthorized, "Bearer token is required")
	}

	if !strings.HasPrefix(tkn, authBearer) {
		return nil, apperr.New(http.StatusUnauthorized, "Bearer token is required")
	}

	did := tkn[len(authBearer)+1:]
	if did == "" {
		return nil, apperr.New(http.StatusUnauthorized, "DID token is required")
	}
	// fmt.Println(did)

	tk, err := token.NewToken(did)
	if err != nil {

		return nil, apperr.New(http.StatusUnauthorized, "Malformed DID token error: "+err.Error())
	}
	// fmt.Println(tk)

	if err := tk.Validate(); err != nil {
		return nil, apperr.New(http.StatusUnauthorized, "DID token failed validation: "+err.Error())
	}

	client := client.New(os.Getenv("MAGIC_API_SECRET"), mag.NewDefaultClient())
	issuer, err := client.User.GetMetadataByIssuer(tk.GetIssuer())
	if err != nil {
		return nil, apperr.New(http.StatusUnauthorized, "Error: "+err.Error())
	}
	// fmt.Println(m, "aaaaaaa")
	// fmt.Println(issuer, "bbbbbbbb")
	if issuer.Email != m.Email {
		return nil, apperr.New(apperr.Unauthorized.Status, "Unauthorized token")
	}

	// find by email
	if user, err := s.userRepo.FindByEmail(issuer.Email); err == nil { // user already exists
		// fmt.Println(user)
		token, expire, err := s.jwt.GenerateToken(user)
		if err != nil {
			return nil, apperr.New(http.StatusUnauthorized, "Unauthorized")
		}
		user.UpdateLastLogin()
		user.Token = xid.New().String()
		if err := s.userRepo.UpdateLogin(user); err != nil {
			return nil, err
		}
		return &model.LoginResponseWithToken{
			Token:        token,
			Expires:      expire,
			RefreshToken: user.Token,
			User:         *user,
		}, nil
	} else {
		u := shortuuid.New()
		user := &model.User{
			Email:        m.Email,
			Verified:     true,
			Active:       true,
			ReferralCode: u,
		}
		userID, err := s.accountRepo.CreateWithMagic(user)
		if err != nil {
			return nil, err
		}
		// fmt.Println(user)
		// fmt.Println(userID)

		// Login
		newUser, err := s.userRepo.View(userID)
		if err == nil { // user already exists
			// fmt.Println(newUser)
			token, expire, err := s.jwt.GenerateToken(newUser)
			if err != nil {
				return nil, apperr.New(http.StatusUnauthorized, "Unauthorized")
			}
			newUser.UpdateLastLogin()
			newUser.Token = xid.New().String()
			if err := s.userRepo.UpdateLogin(newUser); err != nil {
				return nil, err
			}
			return &model.LoginResponseWithToken{
				Token:        token,
				Expires:      expire,
				RefreshToken: newUser.Token,
				User:         *newUser,
			}, nil
		}
	}

	return nil, apperr.NewStatus(http.StatusOK)
}
