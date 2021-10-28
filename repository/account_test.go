package repository_test

import (
	"net/http"
	"testing"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/mock"
	mck "github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/mockgopg"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"

	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AccountUnitTestSuite struct {
	suite.Suite
	mock        *mockgopg.SQLMock
	u           *model.User
	accountRepo *repository.AccountRepo
}

func (suite *AccountUnitTestSuite) SetupTest() {
	var err error
	var db orm.DB
	db, suite.mock, err = mockgopg.NewGoPGDBTest()
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.u = &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "91919191",
	}

	log, _ := zap.NewDevelopment()
	suite.accountRepo = repository.NewAccountRepo(db, log, &mock.Password{})
}

func (suite *AccountUnitTestSuite) TearDownTest() {
	suite.mock.FlushAll()
}

func TestAccountUnitTestSuite(t *testing.T) {
	suite.Run(t, new(AccountUnitTestSuite))
}

// Mock database error when querying
func (suite *AccountUnitTestSuite) TestCreateAndVerifyDBError() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.DB)
	v, err := accountRepo.CreateAndVerify(u)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}

// Mock user already exists
func (suite *AccountUnitTestSuite) TestCreateAndVerifyUserAlreadyExists() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(1, 1, u), nil)
	v, err := accountRepo.CreateAndVerify(u)
	assert.Nil(t, v)
	assert.Equal(t, apperr.New(http.StatusBadRequest, "User already exists."), err)
}

// Mock DB error when inserting user object
func (suite *AccountUnitTestSuite) TestCreateAndVerifyDBErrOnInsertUser() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	mock.ExpectInsert(u).
		Returns(nil, apperr.DB)

	v, err := accountRepo.CreateAndVerify(u)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}

// Mock DB error when inserting verification object
func (suite *AccountUnitTestSuite) TestCreateAndVerifyDBErrOnInsertVerification() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	mock.ExpectInsert(u).
		Returns(nil, nil)

	v := new(model.Verification)
	mock.ExpectInsert(v).
		Returns(nil, apperr.DB)

	v, err := accountRepo.CreateAndVerify(u)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}

// Mock DB error when querying
func (suite *AccountUnitTestSuite) TestCreateWithMobileDBErr() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.DB)
	err := accountRepo.CreateWithMobile(u)
	assert.Equal(t, apperr.DB, err)
}

// Mock user exists and is already verified and active when queried
func (suite *AccountUnitTestSuite) TestCreateWithMobileUserExistsAndVerified() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	u.Verified = true // set expecged user as verified
	u.Active = true   // set expected user object as active
	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(1, 1, u), nil)
	err := accountRepo.CreateWithMobile(u)
	assert.Equal(t, apperr.NewStatus(http.StatusConflict), err)
}

// Mock user exists but is not verified when queried
func (suite *AccountUnitTestSuite) TestCreateWithMobileUserExistsButNotVerified() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(1, 1, nil), nil)
	err := accountRepo.CreateWithMobile(u)
	assert.Equal(t, apperr.BadRequest, err)
}

// Mock HashRandomPassword error
func (suite *AccountUnitTestSuite) TestCreateWithMobileHashRandomPasswordErr() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	// mock a HashRandomPassword error
	accountRepo.Secret = &mck.Password{
		HashRandomPasswordFn: func() (string, error) {
			return "", apperr.DB
		},
	}

	err := accountRepo.CreateWithMobile(u)
	assert.Equal(t, apperr.DB, err)
}

// Mock db error when insert
func (suite *AccountUnitTestSuite) TestCreateWithMobileDBErrOnInsert() {
	u := suite.u
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	// mock a successful HashRandomPassword
	accountRepo.Secret = &mck.Password{
		HashRandomPasswordFn: func() (string, error) {
			return "somerandomhash", nil
		},
	}

	mock.ExpectInsert(u).
		Returns(nil, apperr.DB)

	err := accountRepo.CreateWithMobile(u)
	assert.Equal(t, apperr.DB, err)
}

func (suite *AccountUnitTestSuite) TestFindVerificationTokenSuccess() {
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	var v = new(model.Verification)
	v.Token = "somerandomverificationtoken"
	v.UserID = 1
	mock.ExpectQuery(`SELECT * FROM verifications WHERE (token = ? and deleted_at IS NULL)`).
		WithArgs("somerandomverificationtoken").
		Returns(mockgopg.NewResult(1, 1, v), nil)

	vReturned, err := accountRepo.FindVerificationToken("somerandomverificationtoken")
	assert.Equal(t, v.Token, vReturned.Token)
	assert.Equal(t, v.UserID, vReturned.UserID)
	assert.Nil(t, err)
}

func (suite *AccountUnitTestSuite) TestFindVerificationTokenFailure() {
	accountRepo := suite.accountRepo
	t := suite.T()
	mock := suite.mock

	var v = new(model.Verification)
	v.Token = "anotherverificationtoken"
	v.UserID = 1
	mock.ExpectQuery(`SELECT * FROM verifications WHERE (token = ? and deleted_at IS NULL)`).
		WithArgs("anotherverificationtoken").
		Returns(mockgopg.NewResult(0, 0, v), apperr.NotFound)

	vReturned, err := accountRepo.FindVerificationToken("anotherverificationtoken")
	assert.Nil(t, vReturned)
	assert.Equal(t, apperr.NotFound, err)
}
