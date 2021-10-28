package server

import (
	"os"

	"github.com/alpacahq/ribbit-backend/config"
	"github.com/alpacahq/ribbit-backend/mail"
	mw "github.com/alpacahq/ribbit-backend/middleware"
	"github.com/alpacahq/ribbit-backend/mobile"
	"github.com/alpacahq/ribbit-backend/route"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

// Server holds all the routes and their services
type Server struct {
	RouteServices []route.ServicesI
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Run runs our API server
func (server *Server) Run(env string) error {

	// load configuration
	j := config.LoadJWT(env)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// middleware
	mw.Add(r, CORSMiddleware())
	jwt := mw.NewJWT(j)
	m := mail.NewMail(config.GetMailConfig(), config.GetSiteConfig())
	mobile := mobile.NewMobile(config.GetTwilioConfig())
	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	// setup default routes
	rsDefault := &route.Services{
		DB:     db,
		Log:    log,
		JWT:    jwt,
		Mail:   m,
		Mobile: mobile,
		R:      r}
	rsDefault.SetupV1Routes()

	// setup all custom/user-defined route services
	for _, rs := range server.RouteServices {
		rs.SetupRoutes()
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	// run with port from config
	return r.Run(":" + port)
}
