package postgresDB

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(host, port, user, pass, database, usesSSL string) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, database, usesSSL))
}
