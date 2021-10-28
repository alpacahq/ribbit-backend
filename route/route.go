package route

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/docs"
	"github.com/alpacahq/ribbit-backend/magic"
	"github.com/alpacahq/ribbit-backend/mail"
	mw "github.com/alpacahq/ribbit-backend/middleware"
	"github.com/alpacahq/ribbit-backend/mobile"
	"github.com/alpacahq/ribbit-backend/repository"
	"github.com/alpacahq/ribbit-backend/repository/account"
	assets "github.com/alpacahq/ribbit-backend/repository/assets"
	"github.com/alpacahq/ribbit-backend/repository/auth"
	"github.com/alpacahq/ribbit-backend/repository/plaid"
	"github.com/alpacahq/ribbit-backend/repository/transfer"
	"github.com/alpacahq/ribbit-backend/repository/user"
	"github.com/alpacahq/ribbit-backend/secret"
	"github.com/alpacahq/ribbit-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
	"go.uber.org/zap"
)

// NewServices creates a new router services
func NewServices(DB *pg.DB, Log *zap.Logger, JWT *mw.JWT, Mail mail.Service, Mobile mobile.Service, Magic magic.Service, R *gin.Engine) *Services {
	return &Services{DB, Log, JWT, Mail, Mobile, Magic, R}
}

// Services lets us bind specific services when setting up routes
type Services struct {
	DB     *pg.DB
	Log    *zap.Logger
	JWT    *mw.JWT
	Mail   mail.Service
	Mobile mobile.Service
	Magic  magic.Service
	R      *gin.Engine
}

// SetupV1Routes instances various repos and services and sets up the routers
func (s *Services) SetupV1Routes() {
	// database logic
	userRepo := repository.NewUserRepo(s.DB, s.Log)
	accountRepo := repository.NewAccountRepo(s.DB, s.Log, secret.New())
	assetRepo := repository.NewAssetRepo(s.DB, s.Log, secret.New())
	rbac := repository.NewRBACService(userRepo)

	// s.R.Use(cors.New(cors.Config{
	// 	AllowAllOrigins:  true,
	// 	AllowMethods:     []string{"GET", "PUT", "DELETE", "PATCH", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: false,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return origin == "https://github.com"
	// 	// },
	// 	MaxAge: 12 * time.Hour,
	// }))

	// service logic
	authService := auth.NewAuthService(userRepo, accountRepo, s.JWT, s.Mail, s.Mobile, s.Magic)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac, secret.New())
	userService := user.NewUserService(userRepo, authService, rbac)
	plaidService := plaid.NewPlaidService(userRepo, accountRepo, s.JWT, s.DB, s.Log)
	transferService := transfer.NewTransferService(userRepo, accountRepo, s.JWT, s.DB, s.Log)
	assetsService := assets.NewAssetsService(userRepo, accountRepo, assetRepo, s.JWT, s.DB, s.Log)

	// no prefix, no jwt
	service.AuthRouter(authService, s.R)

	// prefixed with /v1 and protected by jwt
	v1Router := s.R.Group("/v1")
	v1Router.Use(s.JWT.MWFunc())
	service.AccountRouter(accountService, s.DB, v1Router)
	service.PlaidRouter(plaidService, accountService, v1Router)
	service.TransferRouter(transferService, accountService, v1Router)
	service.AssetsRouter(assetsService, accountService, v1Router)
	service.UserRouter(userService, v1Router)

	// Routes for static files
	s.R.StaticFS("/file", http.Dir("public"))
	s.R.StaticFS("/template", http.Dir("templates"))

	//Routes for swagger
	swagger := s.R.Group("swagger")
	{
		docs.SwaggerInfo.Title = "Alpaca MVP"
		docs.SwaggerInfo.Description = "Broker MVP that uses golang gin as webserver, and go-pg library for connecting with a PostgreSQL database"
		docs.SwaggerInfo.Version = "1.0"

		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.R.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
