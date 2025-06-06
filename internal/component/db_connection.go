package component

import (
	"database/sql"
	"fmt"
	"log"

	"adityaad.id/belajar-auth/internal/config"
	_ "github.com/lib/pq"
)

func GetDatabaseConnection(cnf *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s"+
			" port=%s"+
			" user=%s"+
			" password=%s"+
			" dbname=%s"+
			" sslmode=disable",
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.User,
		cnf.Database.Password,
		cnf.Database.Name,
	)

	connection, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err.Error())
	}

	return connection
}
