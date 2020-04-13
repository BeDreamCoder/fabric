/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ormdb

import (
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/config"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb/sqllite3"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

//ORMDBInstance represents a ORMDB instance
type ORMDBInstance struct {
	Config *config.ORMDBConfig
}

//ORMDatabase represents a database within a ORMDB instance
type ORMDatabase struct {
	ORMDBInstance *ORMDBInstance //connection configuration
	DBName        string
	DB            *gorm.DB
	Type          string
}

// NewORMDBInstance create a ORMDB instance through ORMDBConfig
func NewORMDBInstance(metricsProvider metrics.Provider) (*ORMDBInstance, error) {
	ormdbConfig, err := config.GetORMDBConfig()
	if err != nil {

	}
	ormDBInstance := &ORMDBInstance{Config: ormdbConfig}
	return ormDBInstance, nil
}

// CreateORMDatabase creates a ORM database object, as well as the underlying database if it does not exist
func CreateORMDatabase(ormDBInstance *ORMDBInstance, dbName string, dbType string) (*ORMDatabase, error) {
	var db *gorm.DB
	var err error
	switch dbType {
	case "sqlite3":
		db, err = sqllite3.CreateIfNotExistAndOpen(ormDBInstance.Config, dbName)
		if err != nil {
			return nil, errors.WithMessage(err, "create sqlite3 database failed")
		}
	default:
		return nil, errors.New("not supported database type")
	}
	ormDBDatabase := &ORMDatabase{ORMDBInstance: ormDBInstance, DBName: dbName, DB: db, Type: dbType}

	return ormDBDatabase, nil
}

// DeleteORMDatabase deletes a ORM database object, as well as the underlying database if it exists
func DeleteORMDatabase(ormDatabase *ORMDatabase, dbName string) error {
	var err error
	switch ormDatabase.Type {
	case "sqlite3":
		err = sqllite3.DropAndDelete(ormDatabase.DB, ormDatabase.ORMDBInstance.Config, dbName)
		if err != nil {
			return errors.WithMessage(err, "delete sqlite3 database failed")
		}
	default:
		return errors.New("not supported database type")
	}
	return nil
}
