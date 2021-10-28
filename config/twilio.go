package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// TwilioConfig persists the config for our Twilio services
type TwilioConfig struct {
	Account    string `env:"TWILIO_ACCOUNT"`
	Token      string `env:"TWILIO_TOKEN"`
	VerifyName string `env:"TWILIO_VERIFY_NAME"`
	Verify     string `env:"TWILIO_VERIFY"`
}

// GetTwilioConfig returns a TwilioConfig pointer with the correct Mail Config values
func GetTwilioConfig() *TwilioConfig {
	c := TwilioConfig{}

	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	dotenvPath := path.Join(projectRoot, ".env")
	_ = godotenv.Load(dotenvPath)

	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}
