/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package sqllite3

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

// CreateIfNotExistAndOpen opens a gorm sqlite3 database instance
func CreateIfNotExistAndOpen(config *config.ORMDBConfig, dbName string) (db *gorm.DB, err error) {
	_, err = util.CreateDirIfMissing(config.Sqlite3Config.Path)
	if err != nil {
		panic(fmt.Sprintf("error creating sqlite3 db dir if missing: %s", err))
	}
	file := fmt.Sprintf("%s/%s.db", config.Sqlite3Config.Path, dbName)
	db, err = gorm.Open("sqlite3", file)
	return
}

// DropAndDelete deletes a gorm sqlite3 database instance and underlying files
func DropAndDelete(db *gorm.DB, config *config.ORMDBConfig, dbName string) error {
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
