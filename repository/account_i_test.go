package repository_test

import (
	"fmt"
	"log"
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
	"github.com/go-pg/pg/v9/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AccountTestSuite struct {
	suite.Suite
	db       *pg.DB
	dbErr    *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
}

func (suite *AccountTestSuite) SetupTest() {
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
	createSchema(suite.db, &model.Role{}, &model.User{}, &model.Verification{})
}

func (suite *AccountTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func (suite *AccountTestSuite) TestAccountCreateWithMobile() {
	cases := []struct {
		name       string
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.Verification
	}{
		{
			name: "Success: user created",
			user: &model.User{
				CountryCode: "+65",
				Mobile:      "91919191",
			},
			db:        suite.db,
			wantError: nil,
		},
	}
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			accountRepo := repository.NewAccountRepo(tt.db, log, secret.New())
			err := accountRepo.CreateWithMobile(tt.user)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func (suite *AccountTestSuite) TestAccountCreateAndVerify() {
	cases := []struct {
		name       string
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.Verification
	}{
		{
			name: "Success: user created",
			user: &model.User{
				CountryCode: "+65",
				Mobile:      "91919191",
			},
			db:        suite.db,
			wantError: nil,
		},
		{
			name: "Failure: user already exists",
			user: &model.User{
				CountryCode: "+65",
				Mobile:      "91919191",
			},
			db:         suite.db,
			wantError:  apperr.New(http.StatusBadRequest, "User already exists."),
			wantResult: nil,
		},
	}

	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			accountRepo := repository.NewAccountRepo(tt.db, log, secret.New())
			v, err := accountRepo.CreateAndVerify(tt.user)
			assert.Equal(t, tt.wantError, err)
			if v != nil {
				fmt.Println(v.UserID)
				fmt.Println(v.Token)
			}
		})
	}
}

func (suite *AccountTestSuite) TestAccountCreate() {
	cases := []struct {
		name       string
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.User
	}{
		{
			name: "Success: user created",
			user: &model.User{
				Email: "user@example.org",
			},
			db:        suite.db,
			wantError: nil,
			wantResult: &model.User{
				Email: "user@example.org",
			},
		},
		{
			name: "Failure: user already exists",
			user: &model.User{
				Email: "user@example.org",
			},
			db:         suite.db,
			wantError:  apperr.New(http.StatusBadRequest, "User already exists."),
			wantResult: nil,
		},
		{
			name: "Failure: db connection failed",
			db:   suite.dbErr,
			user: &model.User{
				Email: "user2@example.org",
			},
			wantError:  apperr.DB,
			wantResult: nil,
		},
		{
			name: "Failure",
			db:   suite.db,
			user: &model.User{
				ID:    1,
				Email: "user2@example.org",
			},
			wantError:  apperr.DB,
			wantResult: nil,
		},
	}

	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			accountRepo := repository.NewAccountRepo(tt.db, log, secret.New())
			u, err := accountRepo.Create(tt.user)
			assert.Equal(t, tt.wantError, err)
			if u != nil {
				assert.Equal(t, tt.wantResult.Email, u.Email)
			} else {
				assert.Nil(t, u)
			}
		})
	}
}

func (suite *AccountTestSuite) TestChangePasswordSuccess() {
	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(suite.db, log, secret.New())
	userRepo := repository.NewUserRepo(suite.db, log)
	currentPassword := secret.New().HashPassword("currentpassword")
	user := &model.User{
		Email:    "user3@example.org",
		Password: currentPassword,
	}
	u, err := accountRepo.Create(user)
	assert.Equal(suite.T(), user.Password, u.Password)
	assert.NotNil(suite.T(), u.Password)
	assert.Nil(suite.T(), err)

	newPassword := secret.New().HashPassword("newpassword")
	u.Password = newPassword
	assert.NotEqual(suite.T(), currentPassword, newPassword)
	err = accountRepo.ChangePassword(u)
	assert.Nil(suite.T(), err)

	updatedUser, err := userRepo.View(u.ID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), newPassword, updatedUser.Password)

	user2 := &model.User{
		Email:    "user4@example.org",
		Password: currentPassword,
	}
	v, err := accountRepo.CreateAndVerify(user2)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), v)

	vRetrieved, err := accountRepo.FindVerificationToken(v.Token)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), v.Token, vRetrieved.Token)

	err = accountRepo.DeleteVerificationToken(v)
	assert.Nil(suite.T(), err)
}

func (suite *AccountTestSuite) TestChangePasswordFailure() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	accountRepo := repository.NewAccountRepo(suite.dbErr, log, secret.New())
	user := &model.User{
		Email:    "user5@example.org",
		Password: secret.New().HashPassword("somepass"),
	}
	err := accountRepo.ChangePassword(user)
	assert.NotNil(suite.T(), err)
}

func (suite *AccountTestSuite) TestDeleteVerificationTokenFailue() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	accountRepo := repository.NewAccountRepo(suite.dbErr, log, secret.New())
	v := &model.Verification{
		UserID: 1,
		Token:  uuid.NewV4().String(),
	}
	err := accountRepo.DeleteVerificationToken(v)
	assert.NotNil(suite.T(), err)
}

func TestAccountTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
		return
	}
	suite.Run(t, new(AccountTestSuite))
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		opt := &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}
		err := db.CreateTable(model, opt)
		if err != nil {
			log.Fatal(err)
		}
	}
}
