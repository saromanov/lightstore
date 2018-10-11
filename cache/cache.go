package cache

import (
	"errors"

	"github.com/hashicorp/golang-lru"
)

var errNotFound = errors.New("unable to get item")

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

func (cachedata *Cache) Put(item string, timedata int) {
	cachedata.lrudata.Add(item, timedata)
}

func (cachedata *Cache) Get(item string) (interface{}, error) {
	data, ok := cachedata.lrudata.Get(item)
	if !ok {
		return nil, errNotFound
	}
	return data, nil
}
