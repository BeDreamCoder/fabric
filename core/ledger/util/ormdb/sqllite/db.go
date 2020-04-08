package sqllite

import (
	"fmt"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Sqlite3ConnectionDef
type Sqlite3ConnectionDef struct {
	*ormdb.CommonConnectionDef
	Path string
}

// NewSqlite3DB creates a
func NewSqlite3DB(conf *Sqlite3ConnectionDef) (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", fmt.Sprintf("%s/%s", conf.Path, conf.DBName))
	return
}
