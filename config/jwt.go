package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alpacahq/ribbit-backend/secret"

	"github.com/joho/godotenv"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

// LoadJWT returns our JWT with env variables and relevant defaults
func LoadJWT(env string) *JWT {
	jwt := new(JWT)
	defaults.SetDefaults(jwt)

	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	suffix := ""
	if env != "" {
		suffix = suffix + "." + env
	}
	dotenvPath := path.Join(projectRoot, ".env"+suffix)
	_ = godotenv.Load(dotenvPath)

	viper.AutomaticEnv()

	jwt.Secret = viper.GetString("JWT_SECRET")
	if jwt.Secret == "" {
		if strings.HasPrefix(env, "test") {
			// generate jwt secret and write into file
			s, err := secret.GenerateRandomString(256)
			if err != nil {
				log.Fatal(err)
			}
			jwtString := fmt.Sprintf("JWT_SECRET=%s\n", s)
			err = ioutil.WriteFile(dotenvPath, []byte(jwtString), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("Failed to set your environment variable JWT_SECRET. \n" +
				"Please do so via \n" +
				"go run . generate_secret\n" +
				"export JWT_SECRET=[the generated secret]")
		}
	}

	return jwt
}

// JWT holds data necessary for JWT configuration
type JWT struct {
	Realm            string `default:"jwtrealm"`
	Secret           string `default:""`
	Duration         int    `default:"15"`
	RefreshDuration  int    `default:"10"`
	MaxRefresh       int    `default:"10"`
	SigningAlgorithm string `default:"HS256"`
}
