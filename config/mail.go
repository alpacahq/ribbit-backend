package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// MailConfig persists the config for our PostgreSQL database connection
type MailConfig struct {
	Name  string `env:"DEFAULT_NAME"`
	Email string `env:"DEFAULT_EMAIL"`
}

// GetMailConfig returns a MailConfig pointer with the correct Mail Config values
func GetMailConfig() *MailConfig {
	c := MailConfig{}

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
