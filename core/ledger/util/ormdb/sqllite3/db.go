package sqllite3

import (
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"os"
)

// Open opens a gorm sqlite3 database instance
func Open(config *ormdb.ORMDBConfig, dbName string) (db *gorm.DB, err error) {
	dirEmpty, err := util.CreateDirIfMissing(config.Sqlite3Config.Path)
	if err != nil {
		panic(fmt.Sprintf("error creating sqlite3 db dir if missing: %s", err))
	}
	file := fmt.Sprintf("%s/%s.db", config.Sqlite3Config.Path, dbName)
	if dirEmpty {
		_, err = os.Create(file)
		if err != nil {
			panic(fmt.Sprintf("error creating sqlite3 db file: %s", err))
		}
	}
	db, err = gorm.Open("sqlite3", file)
	return
}

func Drop(db *gorm.DB, config *ormdb.ORMDBConfig, dbName string) error {
	err := db.Close()
	if err != nil {
		return errors.WithMessage(err, "error closing sqlite3 db")
	}
	file := fmt.Sprintf("%s/%s.db", config.Sqlite3Config.Path, dbName)
	err = os.RemoveAll(file)
	if err != nil {
		return errors.WithMessage(err, "error removing sqlite3 db file")
	}
	return nil
}
