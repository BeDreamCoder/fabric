package statedb

import (
	"github.com/hyperledger/fabric/common/flogging"
	"sync"
)

var logger = flogging.MustGetLogger("lsccstatecache")

// LsccCacheSize denotes the number of entries allowed in the LsccStateCache
const LsccCacheSize = 50

type LsccStateCache struct {
	Cache   map[string]*VersionedValue
	RWMutex sync.RWMutex
}

func (l *LsccStateCache) GetState(key string) *VersionedValue {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()

	if versionedValue, ok := l.Cache[key]; ok {
		logger.Debugf("key:[%s] found in the LsccStateCache", key)
		return versionedValue
	}
	return nil
}

func (l *LsccStateCache) UpdateState(key string, value *VersionedValue) {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()

	if _, ok := l.Cache[key]; ok {
		logger.Debugf("key:[%s] is updated in LsccStateCache", key)
		l.Cache[key] = value
	}
}

func (l *LsccStateCache) SetState(key string, value *VersionedValue) {
	l.RWMutex.Lock()
	defer l.RWMutex.Unlock()

	if l.IsCacheFull() {
		l.EvictARandomEntry()
	}

	logger.Debugf("key:[%s] is stoed in LsccStateCache", key)
	l.Cache[key] = value
}

func (l *LsccStateCache) IsCacheFull() bool {
	return len(l.Cache) == LsccCacheSize
}

func (l *LsccStateCache) EvictARandomEntry() {
	for key := range l.Cache {
		delete(l.Cache, key)
		return
	}
}
