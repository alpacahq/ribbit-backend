package migration

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alpacahq/ribbit-backend/config"
	"github.com/alpacahq/ribbit-backend/model"

	migrations "github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
  - create_schema [version] - creates initial set of tables from models (structs).
  - sync_assets - sync all the assets from broker.
Usage:
  go run *.go <command> [args]
`

// Run executes migration subcommands
func Run(args ...string) error {
	fmt.Println("Running migration")

	p := config.GetPostgresConfig()

	// connection to db as postgres superuser
	dbSuper := config.GetPostgresSuperUserConnection()
	defer dbSuper.Close()

	// connection to db as POSTGRES_USER
	db := config.GetConnection()
	defer db.Close()

	createUserIfNotExist(dbSuper, p)

	createDatabaseIfNotExist(dbSuper, p)

	if flag.Arg(0) == "create_schema" {
		createSchema(db, &model.Role{}, &model.User{}, &model.Verification{})
		os.Exit(2)
	}

	oldVersion, newVersion, err := migrations.Run(db, args...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
	return nil
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}

func createUserIfNotExist(db *pg.DB, p *config.PostgresConfig) {
	statement := fmt.Sprintf(`SELECT * FROM pg_roles WHERE rolname = '%s';`, p.User)
	res, _ := db.Exec(statement)
	if res.RowsReturned() == 0 {
		statement = fmt.Sprintf(`CREATE USER %s WITH PASSWORD '%s';`, p.User, p.Password)
		_, err := db.Exec(statement)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf(`Created user %s`, p.User)
		}
	}
}

func createDatabaseIfNotExist(db *pg.DB, p *config.PostgresConfig) {
	statement := fmt.Sprintf(`SELECT 1 AS result FROM pg_database WHERE datname = '%s';`, p.Database)
	res, _ := db.Exec(statement)
	if res.RowsReturned() == 0 {
		fmt.Println("creating database")
		statement = fmt.Sprintf(`CREATE DATABASE %s WITH OWNER %s;`, p.Database, p.User)
		_, err := db.Exec(statement)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf(`Created database %s`, p.Database)
		}
	}
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		opt := &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}
		err := db.CreateTable(model, opt)
		if err != nil {
			log.Fatal(err)
		}
	}
}
