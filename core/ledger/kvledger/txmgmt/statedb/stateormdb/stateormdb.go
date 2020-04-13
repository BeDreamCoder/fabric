/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"sync"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"github.com/spf13/viper"
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
	ormDBInstance  *ormdb.ORMDBInstance
	metadataDB     *ormdb.ORMDatabase            // A database per channel to store metadata.
	chainName      string                        // The name of the chain/channel.
	namespaceDBs   map[string]*ormdb.ORMDatabase // One database per deployed chaincode.
	mux            sync.RWMutex
	lsccStateCache *statedb.LsccStateCache
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
	chainName := dbName
	dbName = dbName + "_"

	dbType := viper.GetString("ledger.state.stateDatabase")
	metadataDB, err := ormdb.CreateORMDatabase(ormDBInstance, dbName, dbType)
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*ormdb.ORMDatabase)
	return &VersionedDB{
		ormDBInstance: ormDBInstance,
		metadataDB:    metadataDB,
		chainName:     chainName,
		namespaceDBs:  namespaceDBMap,
		lsccStateCache: &statedb.LsccStateCache{
			Cache: make(map[string]*statedb.VersionedValue),
		},
	}, nil
}

func (v *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	panic("implement me")
}

func (v *VersionedDB) GetVersion(namespace string, key string) (*version.Height, error) {
	panic("implement me")
}

func (v *VersionedDB) GetStateMultipleKeys(namespace string, keys []string) ([]*statedb.VersionedValue, error) {
	panic("implement me")
}

func (v *VersionedDB) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (statedb.ResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) ExecuteQuery(namespace, query string) (statedb.ResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	panic("implement me")
}

func (v *VersionedDB) ApplyUpdates(batch *statedb.UpdateBatch, height *version.Height) error {
	panic("implement me")
}

func (v *VersionedDB) GetLatestSavePoint() (*version.Height, error) {
	panic("implement me")
}

func (v *VersionedDB) ValidateKeyValue(key string, value []byte) error {
	panic("implement me")
}

func (v *VersionedDB) BytesKeySupported() bool {
	panic("implement me")
}

func (v *VersionedDB) Open() error {
	panic("implement me")
}

func (v *VersionedDB) Close() {
	panic("implement me")
}
