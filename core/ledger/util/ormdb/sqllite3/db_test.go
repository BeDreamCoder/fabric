/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sqllite3

import (
	"testing"

	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlite3DB(t *testing.T) {
	conf := &config.ORMDBConfig{Sqlite3Config: &config.Sqlite3Config{Path: "/tmp/ormdb"}}
	db, err := CreateIfNotExistAndOpen(conf, "test")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = DropAndDelete(db, conf, "test")
	assert.NoError(t, err)
}
