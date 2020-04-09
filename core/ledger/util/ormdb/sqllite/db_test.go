package sqllite

import (
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSqlite3DB(t *testing.T) {
	_, err := os.Create("/tmp/test.db")
	conf := &ormdb.ORMDBConfig{DBName: "test.db", Sqlite3Config: &ormdb.Sqlite3Config{Path: "/tmp"}}
	db, err := NewSqlite3DB(conf)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.Close()
	os.RemoveAll("/tmp/test.db")
}
