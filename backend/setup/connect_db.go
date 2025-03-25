package setup

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() (*sql.DB, error) {
	username := os.Getenv("DB_USER_NAME")
	password := os.Getenv("DB_USER_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
