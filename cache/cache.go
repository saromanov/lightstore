package cache

import
(
	"github.com/hashicorp/golang-lru"
)

type Cache struct {
	lrudata *lru
}

func New(size cachesize)*Cache {
	cachedata := new(Cache)
	cachedata.lrudata := lru.New(cachesize)
	return cachedata
}

func (cachedata*Cache) AddToCache(item string,timedata int){
	cachedata.lrudata.Add(item, timedata)
}
