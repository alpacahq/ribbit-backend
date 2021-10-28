package config

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
)

// PostgresConfig persists the config for our PostgreSQL database connection
type PostgresConfig struct {
	URL      string `env:"DATABASE_URL"` // DATABASE_URL will be used in preference if it exists
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

// PostgresSuperUser persists the config for our PostgreSQL superuser
type PostgresSuperUser struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_SUPERUSER" envDefault:"postgres"`
	Password string `env:"POSTGRES_SUPERUSER_PASSWORD" envDefault:""`
	Database string `env:"POSTGRES_SUPERUSER_DB" envDefault:"postgres"`
}

// GetConnection returns our pg database connection
// usage:
// db := config.GetConnection()
// defer db.Close()
func GetConnection() *pg.DB {
	c := GetPostgresConfig()
	// if DATABASE_URL is valid, we will use its constituent values in preference
	validConfig, err := validPostgresURL(c.URL)
	if err == nil {
		c = validConfig
	}
	db := pg.Connect(&pg.Options{
		Addr:     c.Host + ":" + c.Port,
		User:     c.User,
		Password: c.Password,
		Database: c.Database,
		PoolSize: 150,
	})
	return db
}

// GetPostgresConfig returns a PostgresConfig pointer with the correct Postgres Config values
func GetPostgresConfig() *PostgresConfig {
	c := PostgresConfig{}

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

// GetPostgresSuperUserConnection gets the corresponding db connection for our superuser
func GetPostgresSuperUserConnection() *pg.DB {
	c := getPostgresSuperUser()
	db := pg.Connect(&pg.Options{
		Addr:     c.Host + ":" + c.Port,
		User:     c.User,
		Password: c.Password,
		Database: c.Database,
		PoolSize: 150,
	})
	return db
}

func getPostgresSuperUser() *PostgresSuperUser {
	c := PostgresSuperUser{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}

func validPostgresURL(URL string) (*PostgresConfig, error) {
	if URL == "" || strings.TrimSpace(URL) == "" {
		return nil, errors.New("database url is blank")
	}

	validURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	c := &PostgresConfig{}
	c.URL = URL
	c.Host = validURL.Host
	c.Database = validURL.Path
	c.Port = validURL.Port()
	c.User = validURL.User.Username()
	c.Password, _ = validURL.User.Password()
	return c, nil
}
