package e2e_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alpacahq/ribbit-backend/config"
	"github.com/alpacahq/ribbit-backend/e2e"
	"github.com/alpacahq/ribbit-backend/manager"
	mw "github.com/alpacahq/ribbit-backend/middleware"
	"github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/route"
	"github.com/alpacahq/ribbit-backend/secret"
	"github.com/alpacahq/ribbit-backend/testhelper"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var (
	superUser *model.User
	isCI      bool
	port      uint32 = 5432 // uses 5432 in CI, and 9877 when running integration tests locally, against embedded postgresql
)

// end-to-end test constants
const (
	username   string = "db_test_user"
	password   string = "db_test_password"
	database   string = "db_test_database"
	host       string = "localhost"
	tmpDirname string = "tmp2"
)

type E2ETestSuite struct {
	suite.Suite
	db        *pg.DB
	postgres  *embeddedpostgres.EmbeddedPostgres
	m         *manager.Manager
	r         *gin.Engine
	v         *model.Verification
	authToken model.AuthToken
}

// SetupSuite runs before all tests in this test suite
func (suite *E2ETestSuite) SetupSuite() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, tmpDirname)
	os.RemoveAll(tmpDir) // ensure that we start afresh

	_, isCI = os.LookupEnv("CIRCLECI")
	if !isCI { // not in CI environment, so setup our embedded postgresql for integration test
		port = testhelper.AllocatePort(host, 9877)
		testConfig := embeddedpostgres.DefaultConfig().
			Username(username).
			Password(password).
			Database(database).
			Version(embeddedpostgres.V12).
			RuntimePath(tmpDir).
			Port(port)
		suite.postgres = embeddedpostgres.NewDatabase(testConfig)
		err := suite.postgres.Start()
		if err != nil {
			fmt.Println(err)
		}
	}

	suite.db = pg.Connect(&pg.Options{
		Addr:     host + ":" + fmt.Sprint(port),
		User:     username,
		Password: password,
		Database: database,
	})

	log, _ := zap.NewDevelopment()
	defer log.Sync()
	accountRepo := repository.NewAccountRepo(suite.db, log, secret.New())
	roleRepo := repository.NewRoleRepo(suite.db, log)
	suite.m = manager.NewManager(accountRepo, roleRepo, suite.db)

	superUser, _ = e2e.SetupDatabase(suite.m)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// middleware
	mw.Add(r, cors.Default())

	// load configuration
	_ = config.Load("test")
	j := config.LoadJWT("test")
	jwt := mw.NewJWT(j)

	// mock mail
	m := &mock.Mail{
		SendVerificationEmailFn: suite.sendVerification,
	}
	// mock mobile
	mobile := &mock.Mobile{
		GenerateSMSTokenFn: func(string, string) error {
			return nil
		},
		CheckCodeFn: func(string, string, string) error {
			return nil
		},
	}

	// setup routes
	rs := route.NewServices(suite.db, log, jwt, m, mobile, r)
	rs.SetupV1Routes()

	// we can now test our routes in an end-to-end fashion by making http calls
	suite.r = r
}

// TearDownSuite runs after all tests in this test suite
func (suite *E2ETestSuite) TearDownSuite() {
	if !isCI { // not in CI environment, so stop our embedded postgresql db
		suite.postgres.Stop()
	}
}

func (suite *E2ETestSuite) TestGetModels() {
	models := manager.GetModels()
	sql := `SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';`
	var count int
	res, err := suite.db.Query(pg.Scan(&count), sql, nil)

	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(models), count)

	sql = `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`
	var names pg.Strings
	res, err = suite.db.Query(&names, sql, nil)

	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(models), len(names))
}

func (suite *E2ETestSuite) TestSuperUser() {
	assert.NotNil(suite.T(), superUser)
}

func TestE2ETestSuiteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
		return
	}
	suite.Run(t, new(E2ETestSuite))
}

// our mock verification token is saved into suite.token for subsequent use
func (suite *E2ETestSuite) sendVerification(email string, v *model.Verification) error {
	suite.v = v
	return nil
}
