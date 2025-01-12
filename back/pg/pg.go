package pg

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type PG struct {
	Conn            *sql.DB
	QueryMap        map[string]string
	queryMapPrepend string
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
	p := &PG{Conn: db, QueryMap: make(map[string]string, 0)}
	// NOTE: Path is relative to where serve.go is.
	err = p.LoadQueryMap("./pg/queries")
	if err != nil {
		log.Fatalf("could not load query map\n%s", err)
	}
	p.queryMapPrepend = "./pg/queries/"
	return p
}

func (pg *PG) LoadQueryMap(dirPath string) error {
	dirs := []string{dirPath}
	for len(dirs) > 0 {
		dir := pop(&dirs)
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			path := fmt.Sprintf("%s/%s", dir, file.Name())

			if file.IsDir() {
				dirs = append(dirs, path)
				continue
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			extension := strings.Split(file.Name(), ".")
			// Ensure file type (we may want to remove this)
			if extension[len(extension)-1] != "sql" {
				continue
			}

			pg.QueryMap[dir+"/"+extension[0]] = string(data)
		}
	}
	return nil
}

func (pg *PG) QueryRow(queryName string, target any, args ...any) error {
	query, ok := pg.QueryMap[pg.queryMapPrepend+queryName]
	if !ok {
		return fmt.Errorf("Provided query %s/%s does not exist within the query map.", pg.queryMapPrepend, queryName)
	}
	row := pg.Conn.QueryRow(query, args...)
	if err := row.Scan(target); err != nil {
		return fmt.Errorf("Query %s failed with err:\n%s", queryName, err)
	}
	return nil
}

func pop(arr *[]string) string {
	l := len(*arr)
	rv := (*arr)[l-1]
	*arr = (*arr)[:l-1]
	return rv
}
