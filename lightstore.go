package lightstore

import (
	"fmt"
	"github.com/ryszard/goskiplist/skiplist"
	"sync"
	"time"
	"strings"
	//"runtime"
	//"errors"
)

//Basic implmentation of key-value store(without assotiation with any db name)

const (
	param = 0
)

type Settings struct {
	Innerdata string
}

type Store struct {
	items         int
	dbs           map[string]*DB
	mainstore     interface{}
	mainstorename string
	lock          *sync.Mutex
	stat          *Statistics
	index         *Indexing
}

//System keys starts with _ (_index, for example)
func (st*Store) checkSystemKeys(key string) bool {
	return strings.HasPrefix(key, "_")
}

//After understanding, that key is system, make some work with them
func (st*Store) processSystemKey(key string) {
	if key == "_index"{
		st.index.CreateIndex(key)
	}
}

func (st *Store) CheckExistDB(value string) bool {
	_, ok := st.dbs[value]
	return ok
}

func (st *Store) CreateDB(dbname string) {
	st.dbs[dbname] = CreateNewDB(dbname)
}

func (st *Store) Get(value string) interface{} {
	return st.get(value, "")
}

func (st *Store) GetFromDB(dbname string, value string) interface{} {
	return st.get(value, dbname)
}

//if dbname is not equal "", get data from db with name dbname
func (st *Store) get(value string, dbname string) interface{} {
	mainstore := st.mainstore
	if dbname != "" {
		//check db availability
		dbdata, ok := st.dbs[dbname]
		if !ok {
			log.Fatal(fmt.Sprintf("db with name %s is not found", dbname))
		} else if !dbdata.isactive {
			log.Fatal(fmt.Sprintf("db with name %s is not active", dbname))
		} else {
			mainstore = dbdata.mainstore
		}
	}

	st.lock.Lock()
	defer st.lock.Unlock()
	switch mainstore.(type) {
	case *Dict:
		result, ok := mainstore.(*Dict).Get(value)
		if ok {
			st.stat.num_reads += 1
			return result
		}
	case *BMtree:
		fmt.Println("Not implemented yet")
	case *skiplist.SkipList:
		result, ok := mainstore.(*skiplist.SkipList).Get(value)
		if ok {
			st.stat.num_reads += 1
			return result
		}
	}

	return nil

}

//Get many kayes from list
func (st *Store) GetMany(keys []string) interface{} {
	result := make([]interface{}, len(keys))
	st.lock.Lock()
	defer st.lock.Unlock()
	if len(keys) > 0 {
		for i := 0; i < len(keys); i++ {
			result = append(result, st.Get(keys[i]))
		}

		return result
	}
	return nil
}

func (st *Store) Set(key string, value interface{}) bool {
	return st.set("", key, value)
}

func (st *Store) SetinDB(dbname string, key string, value interface{}) bool {
	return st.set(dbname, key, value)
}

func (st *Store) set(dbname string, key string, value interface{}) bool {
	st.lock.Lock()
	defer st.lock.Unlock()
	if st.checkSystemKeys(key) {
		return true
	}
	mainstore := st.mainstore
	if dbname != "" {
		//check db availability
		dbdata, ok := st.dbs[dbname]
		if !ok {
			log.Info(fmt.Sprintf("db with name %s is not found", dbname))
			return false
		} else if !dbdata.isactive {
			log.Info(fmt.Sprintf("db with name %s is not active", dbname))
			return false
		} else if dbdata.limit != -1 && (dbdata.limit-dbdata.datacount) == 0 {
			log.Info(fmt.Sprintf("db with name %s not availability to write, because contains maxumum possible number of data objrct"))
			return false
		} else {
			mainstore = dbdata.mainstore
		}
	}

	switch mainstore.(type) {
	case *Dict:
		mainstore.(*Dict).Set(key, value)
	case *skiplist.SkipList:
		mainstore.(*skiplist.SkipList).Set(key, value)
	}

	if dbname != "" {
		dbdata, _ := st.dbs[dbname]
		dbdata.datacount += 1
	}
	st.stat.num_writes += 1
	return true

}

func (st *Store) Remove(key string) {
	switch st.mainstore.(type) {
	case Dict:
		st.mainstore.(*Dict).Remove(key)
	default:
		st.mainstore.(*skiplist.SkipList).Delete(key)
	}
}

func (st *Store) Find(key string) interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	switch st.mainstore.(type) {

	}

	return nil
}

//Return statistics of usage
func (st *Store) Stat() *Statistics {
	return st.stat
}

func (st *Store) CloseLightStore() {
	fmt.Println("End working: ", time.Now())
}

//This private method provides checking inner datastructure for storing
func checkDS(name string) (result interface{}) {
	result = skiplist.NewStringMap()
	/* SkipList datastructure as main store */
	if name == "skiplist" {
		store.mainstore = skiplist.NewStringMap()
	}

	/* Simple map as main store */
	if name == "dict" {
		store.mainstore = NewDict()
	}

	/*B-tree structure as main store */
	if name == "b-tree" {
		store.mainstore = InitBMTree()
	}
	return result
}

func InitStore(settings Settings) *Store {
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	mutex := &sync.Mutex{}
	store := new(Store)
	starttime := time.Now()
	fmt.Println("Lightstore is working: ", starttime)
	store.items = 0
	store.mainstore = checkDS(settings.Innerdata)
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(Statistics)
	store.stat.start = starttime
	store.index = NewIndexing()
	return store
}
