/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"github.com/hyperledger/fabric/core/ledger/util/ormdb"
)

// VersionedDBProvider implements interface VersionedDBProvider
type VersionedDBProvider struct {
	ormDBinstance *ormdb.ORMDBInstance
	databases     map[string]*VersionedDB
}

// VersionedDB implements VersionedDB interface
type VersionedDB struct {
}
