package vm

import (
	"sync"

	"github.com/Giulio2002/gevm/types"
)

const globalJumpDestCacheLimit = 32768

var globalJumpDestCache = struct {
	sync.RWMutex
	items map[types.B256][]byte
	order []types.B256
	next  int
}{
	items: make(map[types.B256][]byte, globalJumpDestCacheLimit),
	order: make([]types.B256, 0, globalJumpDestCacheLimit),
}

func GetCachedJumpDest(hash types.B256) ([]byte, bool) {
	globalJumpDestCache.RLock()
	jt, ok := globalJumpDestCache.items[hash]
	globalJumpDestCache.RUnlock()
	return jt, ok
}

func PutCachedJumpDest(hash types.B256, jt []byte) {
	if len(jt) == 0 {
		return
	}
	globalJumpDestCache.Lock()
	defer globalJumpDestCache.Unlock()
	if _, ok := globalJumpDestCache.items[hash]; ok {
		globalJumpDestCache.items[hash] = jt
		return
	}
	if len(globalJumpDestCache.order) >= globalJumpDestCacheLimit {
		evict := globalJumpDestCache.order[globalJumpDestCache.next]
		globalJumpDestCache.order[globalJumpDestCache.next] = hash
		delete(globalJumpDestCache.items, evict)
		globalJumpDestCache.next++
		if globalJumpDestCache.next == globalJumpDestCacheLimit {
			globalJumpDestCache.next = 0
		}
	} else {
		globalJumpDestCache.order = append(globalJumpDestCache.order, hash)
	}
	globalJumpDestCache.items[hash] = jt
}
