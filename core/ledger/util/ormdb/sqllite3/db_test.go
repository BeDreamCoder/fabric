package sqllite3

import (
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSqlite3DB(t *testing.T) {
	conf := &ormdb.ORMDBConfig{Sqlite3Config: &ormdb.Sqlite3Config{Path: "/tmp"}}
	db, err := Open(conf, "test")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	_ = Drop(db, conf, "test")
}
