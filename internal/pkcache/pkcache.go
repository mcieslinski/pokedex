package pkcache

import (
  "fmt"
  "sync"
  "time"
)

type pkCacheEntry struct {
  createdAt time.Time
  bytes     []byte
}

type PkCache struct {
  cache  map[string]pkCacheEntry
  mux    sync.RWMutex
  tickr  *time.Ticker
  expiry time.Duration
}

func (pc *PkCache) Add(key string, value []byte) {
  pc.mux.Lock()
  defer pc.mux.Unlock()

  pc.cache[key] = pkCacheEntry{ createdAt: time.Now(), bytes: value, }
}

func (pc *PkCache) Get(key string) ([]byte, bool) {
  pc.mux.RLock()
  defer pc.mux.RUnlock()

  res, ok := pc.cache[key]
  if ok { fmt.Println("Cache hit on:", key) }
  return res.bytes, ok
}

func (pc *PkCache) reapLoop() {
  for {
    now := <- pc.tickr.C

    pc.mux.Lock()
    for key, item := range pc.cache {
      if now.After(item.createdAt.Add(pc.expiry)) {
        delete(pc.cache, key)
      }
    }
    pc.mux.Unlock()
  }
}

func NewPkCache(interval time.Duration) *PkCache {
  ret := &PkCache{
    cache:  map[string]pkCacheEntry{},
    mux:    sync.RWMutex{},
    tickr:  time.NewTicker(100 * time.Millisecond),
    expiry: interval,
  }

  go ret.reapLoop()

  return ret
}
