package pg

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PG struct {
	Conn *sql.DB
}

func NewPG() *PG {
	connStr := os.Getenv("PG_CONN_STRING")
	if connStr == "" {
		log.Fatal("could not read PG_CONN_STRING from .env\n")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not open connection to db\n%s", err)
	}
	return &PG{Conn: db}
}
