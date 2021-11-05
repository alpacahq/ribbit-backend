package repository_test

import (
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/repository/account"
	"github.com/alpacahq/ribbit-backend/secret"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type RBACTestSuite struct {
	suite.Suite
	db       *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
}

func (suite *RBACTestSuite) SetupTest() {
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
	createSchema(suite.db, &model.Role{}, &model.User{}, &model.Verification{})
}

func (suite *RBACTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func TestRBACTestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
		return
	}
	suite.Run(t, new(RBACTestSuite))
}

func (suite *RBACTestSuite) TestRBAC() {
	// create a context for tests
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Set("role", int8(model.SuperAdminRole))

	// create a user in our test database, which is superadmin
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.db, log)
	accountRepo := repository.NewAccountRepo(suite.db, log, secret.New())
	rbac := repository.NewRBACService(userRepo)

	// ensure that our roles table is populated with default roles
	roleRepo := repository.NewRoleRepo(suite.db, log)
	err := roleRepo.CreateRoles()
	assert.Nil(suite.T(), err)

	accountService := account.NewAccountService(userRepo, accountRepo, rbac, secret.New())
	err = accountService.Create(c, &model.User{
		CountryCode: "+65",
		Mobile:      "91919191",
		Active:      true,
		RoleID:      3,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), rbac)

	// since the current user is a superadmin, we should be able to change user data
	userID := 1
	access := rbac.EnforceUser(c, userID)
	assert.True(suite.T(), access)

	// since the current user is a superadmin, we should be able to change location data
	// access = rbac.EnforceLocation(c, 1)
	// assert.True(suite.T(), access)
}
