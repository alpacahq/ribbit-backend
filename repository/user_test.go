package repository_test

import (
	"testing"

	"github.com/alpacahq/ribbit-backend/mockgopg"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"

	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserUnitTestSuite struct {
	suite.Suite
	mock     *mockgopg.SQLMock
	u        *model.User
	userRepo *repository.UserRepo
}

func (suite *UserUnitTestSuite) SetupTest() {
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
		Token:       "someusertoken",
	}

	log, _ := zap.NewDevelopment()
	suite.userRepo = repository.NewUserRepo(db, log)
}

func (suite *UserUnitTestSuite) TearDownTest() {
	suite.mock.FlushAll()
}

func TestUserUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UserUnitTestSuite))
}

func (suite *UserUnitTestSuite) TestFindByReferralCodeSuccess() {
	u := suite.u
	userRepo := suite.userRepo
	t := suite.T()
	mock := suite.mock

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	mock.ExpectQueryOne(sql).
		WithArgs(u.Username).
		Returns(mockgopg.NewResult(1, 1, u), nil)

	uReturned, err := userRepo.FindByReferralCode("hello")
	assert.Equal(t, u.Username, uReturned.Username)
	assert.Nil(t, err)
}

func (suite *UserUnitTestSuite) TestFindByUsernameSuccess() {
	u := suite.u
	userRepo := suite.userRepo
	t := suite.T()
	mock := suite.mock

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	mock.ExpectQueryOne(sql).
		WithArgs(u.Username).
		Returns(mockgopg.NewResult(1, 1, u), nil)

	uReturned, err := userRepo.FindByUsername("hello")
	assert.Equal(t, u.Username, uReturned.Username)
	assert.Nil(t, err)
}

func (suite *UserUnitTestSuite) TestFindByEmailSuccess() {
	u := suite.u
	userRepo := suite.userRepo
	t := suite.T()
	mock := suite.mock

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."email" = ? and deleted_at is null)`
	mock.ExpectQueryOne(sql).
		WithArgs(u.Email).
		Returns(mockgopg.NewResult(1, 1, u), nil)

	uReturned, err := userRepo.FindByEmail("hello@world.org")
	assert.Equal(t, u.Email, uReturned.Email)
	assert.Nil(t, err)
}

func (suite *UserUnitTestSuite) TestFindByMobileSuccess() {
	u := suite.u
	userRepo := suite.userRepo
	t := suite.T()
	mock := suite.mock

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
        FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
        WHERE ("user"."country_code" = ? and "user"."mobile" = ? and deleted_at is null)`
	mock.ExpectQueryOne(sql).
		WithArgs(u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(1, 1, u), nil)

	uReturned, err := userRepo.FindByMobile(u.CountryCode, u.Mobile)
	assert.Equal(t, u.Mobile, uReturned.Mobile)
	assert.Nil(t, err)
}

func (suite *UserUnitTestSuite) TestFindByTokenSuccess() {
	u := suite.u
	userRepo := suite.userRepo
	t := suite.T()
	mock := suite.mock

	u.Token = "someusertoken"

	var user = new(model.User)
	user.Token = "someusertoken"
	user.ID = 1
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	mock.ExpectQueryOne(sql).
		WithArgs("someusertoken").
		Returns(mockgopg.NewResult(1, 1, user), nil)

	_, err := userRepo.FindByToken(u.Token)
	assert.Equal(t, u.Token, user.Token)
	assert.Nil(t, err)
}
