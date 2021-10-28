package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alpacahq/ribbit-backend/config"
	mw "github.com/alpacahq/ribbit-backend/middleware"
	"github.com/alpacahq/ribbit-backend/mock"
	"github.com/alpacahq/ribbit-backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func hwHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"text": "Hello World.",
	})
}

func ginHandler(mw ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	for _, v := range mw {
		r.Use(v)
	}
	r.GET("/hello", hwHandler)
	return r
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		header     string
	}{
		{
			name:       "Empty header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header not containing Bearer",
			header:     "notBearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid header",
			header:     mock.HeaderInvalid(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Success",
			header:     mock.HeaderValid(),
			wantStatus: http.StatusOK,
		},
	}
	jwtCfg := &config.JWT{Realm: "testRealm", Secret: "jwtsecret", Duration: 60, SigningAlgorithm: "HS256"}
	jwtMW := mw.NewJWT(jwtCfg)
	ts := httptest.NewServer(ginHandler(jwtMW.MWFunc()))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		name      string
		wantToken string
		req       *model.User
	}{
		{
			name: "Success",
			req: &model.User{
				Base:     model.Base{},
				ID:       1,
				Username: "johndoe",
				Email:    "johndoe@mail.com",
				Role: &model.Role{
					AccessLevel: model.SuperAdminRole,
				},
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}
	jwtCfg := &config.JWT{Realm: "testRealm", Secret: "jwtsecret", Duration: 60, SigningAlgorithm: "HS256"}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			jwt := mw.NewJWT(jwtCfg)
			str, _, err := jwt.GenerateToken(tt.req)
			assert.Nil(t, err)
			assert.Equal(t, tt.wantToken, strings.Split(str, ".")[0])
		})
	}
}
