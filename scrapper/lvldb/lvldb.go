// In the case of the scrapper all this is meant to do is
// keep a persistant k/v storage of images we have uploaded
// so we dont run into duplicates.
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

func (db *LvlDB) Exists(key string) (bool, error) {
	_, err := db.Conn.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (db *LvlDB) Create(key string) error {
	err := db.Conn.Put([]byte(key), []byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}
