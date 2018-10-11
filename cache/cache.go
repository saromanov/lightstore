package cache

import "github.com/hashicorp/golang-lru"

type Cache struct {
	lrudata *lru.Cache
}

func New(cachesize int) *Cache {
	cachedata := new(Cache)
	cachedata.lrudata, _ = lru.New(cachesize)
	return cachedata
}

func (cachedata *Cache) AddToCache(item string, timedata int) {
	cachedata.lrudata.Add(item, timedata)
}
