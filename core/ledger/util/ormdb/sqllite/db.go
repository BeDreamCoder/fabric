package sqllite

import (
	"fmt"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)


// NewSqlite3DB creates a
func NewSqlite3DB(config *ormdb.ORMDBConfig) (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", fmt.Sprintf("%s/%s", config.Sqlite3Config.Path, config.DBName))
	return
}
