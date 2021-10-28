package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// MagicConfig persists the config for our Magic services
type MagicConfig struct {
	Key    string `env:"MAGIC_API_KEY"`
	Secret string `env:"MAGIC_API_SECRET"`
}

// GetMagicConfig returns a MagicConfig pointer with the correct Magic.link Config values
func GetMagicConfig() *MagicConfig {
	c := MagicConfig{}

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
