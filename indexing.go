package lightstore

import (
	"./cache"
)

//Indexing in lightstore

type Indexing struct {
	index      map[string][]string
	location   map[string]string
	caching    *cache.Cache
	maxsize    int
	maxindexes int
}

func NewIndexing() *Indexing {
	idx := new(Indexing)
	idx.index = make(map[string][]string)
	idx.location = make(map[string]string)
	idx.caching = cache.New(1000)
	idx.maxsize = 1000000
	idx.maxindexes = 10000
	return idx
}

func (idx *Indexing) CreateIndex(name string) {
	idx.index[name] = []string{}
}

func (idx* Indexing) DropIndex(value string) {
	idx.DropIndexes([]string{value})
}

func (idx*Indexing) DropIndexes(values []string) {
	for key, _ := range idx.index {
		for _, value := range values {
			if key == value {
				delete(idx.index, key)
			}
		}
	}
}

func (idx *Indexing) AddItem(name, value, location string) {
	idx.index[name] = append(idx.index[name], value)
	idx.location[value] = location
}

func (idx *Indexing) AddItemToCache(name, value, location string) {
	idx.caching.AddToCache(value, 1000)
}
