/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"sync"
)

var logger = flogging.MustGetLogger("stateormdb")

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBinstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	//openCounts    uint64
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	ormDBInstance *ormdb.ORMDBInstance
	metadataDB    *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName     string                        // The name of the chain/channel.
	namespaceDBs  map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	rwmutex       sync.RWMutex
}

// NewVersionedDBProvider instantiates VersionedDBProvider
func NewVersionedDBProvider(metricsProvider metrics.Provider) (*VersionedDBProvider, error) {
	logger.Debugf("constructing ORMDB VersionedDBProvider")
	instance, err := ormdb.NewORMDBInstance(metricsProvider)
	if err != nil {
		return nil, err
	}
	return &VersionedDBProvider{ormDBinstance: instance, databases: make(map[string]*VersionedDB)}, nil
}

// GetDBHandle gets the handle to a named database
func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()
	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(provider.ormDBinstance, dbName)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}

// newVersionedDB constructs an instance of VersionedDB
func newVersionedDB(ormDBInstance *ormdb.ORMDBInstance, dbName string) (*VersionedDB, error) {
	// CreateCouchDatabase creates a CouchDB database object, as well as the underlying database if it does not exist
	chainName := dbName
	dbName = dbName + "_"

	metadataDB, err := ormdb.NewORMDBInstance()CreateCouchDatabase(ormDBInstance, dbName)
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*couchdb.CouchDatabase)
	return &VersionedDB{
		couchInstance:      couchInstance,
		metadataDB:         metadataDB,
		chainName:          chainName,
		namespaceDBs:       namespaceDBMap,
		committedDataCache: newVersionCache(),
		lsccStateCache: &lsccStateCache{
			cache: make(map[string]*statedb.VersionedValue),
		},
	}, nil
}
