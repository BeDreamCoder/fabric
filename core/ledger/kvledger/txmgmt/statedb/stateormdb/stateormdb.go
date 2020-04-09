/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
	"sync"
)

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBinstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
	mux           sync.Mutex
	openCounts    uint64
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
	couchInstance      *couchdb.CouchInstance
	metadataDB         *couchdb.CouchDatabase            // A database per channel to store metadata such as savepoint.
	chainName          string                            // The name of the chain/channel.
	namespaceDBs       map[string]*couchdb.CouchDatabase // One database per deployed chaincode.
	committedDataCache *versionsCache                    // Used as a local cache during bulk processing of a block.
	verCacheLock       sync.RWMutex
	mux                sync.RWMutex
	lsccStateCache     *lsccStateCache
}
