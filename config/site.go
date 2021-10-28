package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// SiteConfig persists global configs needed for our application
type SiteConfig struct {
	ExternalURL string `env:"EXTERNAL_URL"  envDefault:"http://localhost:8080"`
}

// GetSiteConfig returns a SiteConfig pointer with the correct Site Config values
func GetSiteConfig() *SiteConfig {
	c := SiteConfig{}

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
