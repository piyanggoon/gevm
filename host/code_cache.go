package host

import (
	"sync"

	"github.com/Giulio2002/gevm/types"
	"github.com/Giulio2002/gevm/vm"
)

const globalCodeCacheLimit = 32768

type cachedCode struct {
	padded types.Bytes
	length int
}

var globalCodeCache = struct {
	sync.RWMutex
	items map[types.B256]cachedCode
	order []types.B256
	next  int
}{
	items: make(map[types.B256]cachedCode, globalCodeCacheLimit),
	order: make([]types.B256, 0, globalCodeCacheLimit),
}

func getCachedCode(hash types.B256) (types.Bytes, bool) {
	globalCodeCache.RLock()
	item, ok := globalCodeCache.items[hash]
	globalCodeCache.RUnlock()
	if !ok {
		return nil, false
	}
	return item.padded[:item.length], true
}

func getCachedPaddedCode(hash types.B256) (types.Bytes, int, bool) {
	globalCodeCache.RLock()
	item, ok := globalCodeCache.items[hash]
	globalCodeCache.RUnlock()
	return item.padded, item.length, ok
}

func putCachedCode(hash types.B256, code types.Bytes) {
	if hash == types.B256Zero || hash == types.KeccakEmpty || len(code) == 0 {
		return
	}
	padded := make(types.Bytes, len(code)+vm.BytecodeEndPadding)
	copy(padded, code)
	item := cachedCode{padded: padded, length: len(code)}

	globalCodeCache.Lock()
	defer globalCodeCache.Unlock()
	if _, ok := globalCodeCache.items[hash]; ok {
		globalCodeCache.items[hash] = item
		return
	}
	if len(globalCodeCache.order) >= globalCodeCacheLimit {
		evict := globalCodeCache.order[globalCodeCache.next]
		globalCodeCache.order[globalCodeCache.next] = hash
		delete(globalCodeCache.items, evict)
		globalCodeCache.next++
		if globalCodeCache.next == globalCodeCacheLimit {
			globalCodeCache.next = 0
		}
	} else {
		globalCodeCache.order = append(globalCodeCache.order, hash)
	}
	globalCodeCache.items[hash] = item
}
