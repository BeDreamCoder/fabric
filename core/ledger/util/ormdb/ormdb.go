package ormdb

import (
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/jinzhu/gorm"
)

//ORMDBInstance represents a ORMDB instance
type ORMDBInstance struct {
	Config *ORMDBConfig
}

//CouchDatabase represents a database within a CouchDB instance
type ORMDatabase struct {
	ORMDBInstance *ORMDBInstance //connection configuration
	DBName        string
	DB            *gorm.DB
	Type          string
}

//NewORMDBInstance create a ORMDB instance through ORMDBConfig
func NewORMDBInstance(metricsProvider metrics.Provider) (*ORMDBInstance, error) {
	config, err := GetORMDBConfig()
	if err != nil {

	}
	ormDBInstance := &ORMDBInstance{Config: config}
	return ormDBInstance, nil
}

//CreateCouchDatabase creates a CouchDB database object, as well as the underlying database if it does not exist
func CreateCouchDatabase(couchInstance *CouchInstance, dbName string) (*CouchDatabase, error) {

	databaseName, err := mapAndValidateDatabaseName(dbName)
	if err != nil {
		logger.Errorf("Error calling CouchDB CreateDatabaseIfNotExist() for dbName: %s, error: %s", dbName, err)
		return nil, err
	}

	couchDBDatabase := CouchDatabase{CouchInstance: couchInstance, DBName: databaseName, IndexWarmCounter: 1}

	// Create CouchDB database upon ledger startup, if it doesn't already exist
	err = couchDBDatabase.CreateDatabaseIfNotExist()
	if err != nil {
		logger.Errorf("Error calling CouchDB CreateDatabaseIfNotExist() for dbName: %s, error: %s", dbName, err)
		return nil, err
	}

	return &couchDBDatabase, nil
}
