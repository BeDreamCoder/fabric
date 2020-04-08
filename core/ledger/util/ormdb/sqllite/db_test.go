package sqllite

import (
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSqlite3DB(t *testing.T) {
	conf := &Sqlite3ConnectionDef{&ormdb.CommonConnectionDef{DBName: "test"},
		"test"}
	db, err := NewSqlite3DB(conf)
	assert.NoError(t, err)
	db.Close()
}
