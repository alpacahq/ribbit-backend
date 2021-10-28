package manager

import (
	"fmt"

	"github.com/alpacahq/ribbit-backend/config"

	"github.com/go-pg/pg/v9"
)

// CreateDatabaseIfNotExist creates our postgresql database from postgres config
func CreateDatabaseIfNotExist(db *pg.DB, p *config.PostgresConfig) {
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
	} else {
		fmt.Printf("Database named %s already exists\n", p.Database)
	}
}
