// This is basically just a persistant cache to track what images have already been
// downloaded and what taskIdUrls have been fetched so we can run this as a cron
// script and get newer data from them.
package lvldb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type LvlDB struct {
	Conn *leveldb.DB
}

func NewLvlDB() LvlDB {
	db, err := leveldb.OpenFile("lvldb/db", nil)
	if err != nil {
		panic(fmt.Sprintf("Could not open leveldb\n%s", err))
	}
	return LvlDB{
		Conn: db,
	}
}

// simple get wrapper for lvldb
func (db *LvlDB) Get(key string) (string, error) {
	timestamp, err := db.Conn.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(timestamp), nil
}

// simple put wrapper for lvldb
func (db *LvlDB) Put(key string, value string) error {
	err := db.Conn.Put([]byte(key), []byte(value), nil)
	// TODO: Err handling
	if err != nil {
		return err
	}
	return nil
}
