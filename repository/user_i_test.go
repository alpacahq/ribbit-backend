package repository_test

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/secret"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserTestSuite struct {
	suite.Suite
	db       *pg.DB
	dbErr    *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
	u        *model.User // test user
}

func (suite *UserTestSuite) SetupTest() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, "tmp")
	os.RemoveAll(tmpDir)
	testConfig := embeddedpostgres.DefaultConfig().
		Username("db_test_user").
		Password("db_test_password").
		Database("db_test_database").
		Version(embeddedpostgres.V12).
		RuntimePath(tmpDir).
		Port(9876)

	suite.postgres = embeddedpostgres.NewDatabase(testConfig)
	err := suite.postgres.Start()
	assert.Equal(suite.T(), err, nil)

	suite.db = pg.Connect(&pg.Options{
		Addr:     "localhost:9876",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})
	suite.dbErr = pg.Connect(&pg.Options{
		Addr:     "localhost:9875",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})
	suite.u = &model.User{
		Username:    "user",
		Email:       "user@example.org",
		CountryCode: "+65",
		Mobile:      "91919191",
	}
	createSchema(suite.db, &model.Role{}, &model.User{}, &model.Verification{})
}

func (suite *UserTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func (suite *UserTestSuite) TestUserView() {
	cases := []struct {
		name       string
		create     bool
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.Verification
	}{
		{
			name:      "Fail: user not found",
			create:    false,
			user:      suite.u,
			db:        suite.db,
			wantError: &apperr.APPError{Status: http.StatusNotFound, Message: "400 not found"},
		},
		{
			name:      "Success: view user, find user",
			create:    true,
			user:      suite.u,
			db:        suite.db,
			wantError: nil,
		},
	}
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			userRepo := repository.NewUserRepo(tt.db, log)

			if tt.create {
				accountRepo := repository.NewAccountRepo(tt.db, log, secret.New())
				_, err := accountRepo.Create(tt.user)
				assert.Nil(t, err)
				u, err := userRepo.View(tt.user.ID)
				assert.Nil(t, err)
				assert.Equal(t, tt.user.Mobile, u.Mobile)
				assert.False(t, u.Active)
				assert.False(t, u.Verified)
				assert.Nil(t, u.LastLogin)
				err = userRepo.UpdateLogin(u)
				assert.Nil(t, err)
				u, err = userRepo.View(u.ID)
				assert.NotNil(t, u.LastLogin)
				u.Active = true
				u.Verified = true
				u, err = userRepo.Update(u)
				assert.Nil(t, err)
				u, err = userRepo.View(u.ID)
				assert.Nil(t, err)
				assert.True(t, u.Active)
				assert.True(t, u.Verified)
				err = userRepo.Delete(u)
				assert.Nil(t, err)
				u, err = userRepo.View(u.ID)
				assert.Nil(t, u)
				assert.Error(t, apperr.NotFound)
				pag := &model.Pagination{Limit: 10, Offset: 0}
				users, err := userRepo.List(nil, pag)
				assert.Equal(suite.T(), 0, len(users))
				assert.Nil(suite.T(), err)
			} else {
				u, err := userRepo.View(tt.user.ID)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByMobile(tt.user.CountryCode, tt.user.Mobile)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByEmail(tt.user.Email)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				findByReferralCodeResponse, err := userRepo.FindByReferralCode(tt.user.ReferralCode)
				assert.Nil(t, findByReferralCodeResponse)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByUsername(tt.user.Username)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByToken("somerandomtokenthatdoesntexist")
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
			}
		})
	}
}

func (suite *UserTestSuite) TestUpdateLoginFailure() {
	u := suite.u
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.dbErr, log)
	err := userRepo.UpdateLogin(u)
	assert.NotNil(suite.T(), err)
}

func (suite *UserTestSuite) TestUpdateFailure() {
	u := suite.u
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.dbErr, log)
	u.Address = "some address"
	user, err := userRepo.Update(u)
	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserTestSuite) TestDeleteFailure() {
	u := suite.u
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.dbErr, log)
	err := userRepo.Delete(u)
	assert.NotNil(suite.T(), err)
}

func (suite *UserTestSuite) TestListFailure() {
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.dbErr, log)
	qp := &model.ListQuery{}
	pag := &model.Pagination{Limit: 10, Offset: 0}
	_, err := userRepo.List(qp, pag)
	assert.NotNil(suite.T(), err)
}

func TestUserTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
		return
	}
	suite.Run(t, new(UserTestSuite))
}
