package cache

import "github.com/hashicorp/golang-lru"

type Cache struct {
	lrudata   *lru.Cache
	cachesize int
}

func New(cachesize int) *Cache {
	cachedata := new(Cache)
	cachedata.cachesize = cachesize
	cachedata.lrudata, _ = lru.New(cachesize)
	return cachedata
}

func (cachedata *Cache) AddToCache(item string, timedata int) {
	cachedata.lrudata.Add(item, timedata)
}
