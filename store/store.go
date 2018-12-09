package store

import (
	"fmt"
	"sync"
	"time"

	ds "github.com/saromanov/lightstore/datastructures"
	"github.com/saromanov/lightstore/history"
	log "github.com/saromanov/lightstore/logging"
	"github.com/saromanov/lightstore/monitoring"
	"github.com/saromanov/lightstore/rpc"
	"github.com/saromanov/lightstore/statistics"
)

//Basic implmentation of key-value store(without assotiation with any db name)

const param = 0

const (
	MaxKeySize   = 512
	MaxValueSize = 32768
)

type Settings struct {
	Innerdata string
}

// Store provides implementation of the main store
type Store struct {
	items     int
	dbs       map[string]*DB
	store     ds.Storage
	storename string
	keys      []string
	lock      *sync.RWMutex
	stat      *statistics.Statistics
	index     *Indexing
	config    *Config
	pubsub    *Pubsub
	//Event history
	historyevent *history.History
	rpcdata      *rpc.RPCData
	compression  bool
	fileWatcher  *watcher
}

// newStore creates a new instance of lightstore
func newStore(c *Config) *Store {
	if c == nil {
		c = LoadConfigData("")
	}
	mutex := &sync.RWMutex{}
	store := new(Store)
	starttime := time.Now().UTC()
	store.items = 0
	fileWatcher, err := newWatcher(".")
	if err != nil {
		log.Info(fmt.Sprintf("unable to init file watcher: %v", err))
	}
	if c.Monitoring {
		monitoring.Init()
	}
	store.fileWatcher = fileWatcher
	store.store = checkDS("")
	store.keys = []string{}
	store.compression = c.Compression
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(statistics.Statistics)
	store.stat.Start = starttime
	store.index = NewIndexing()
	store.config = LoadConfigData("")
	store.historyevent = history.NewHistory(5)
	store.pubsub = PubsubInit()
	rpc.RegisterRPCFunction(store.pubsub)
	store.rpcdata = rpc.Init("")
	store.rpcdata.Run()
	return store
}

//After understanding, that key is system, make some work with them
func (st *Store) processSystemKey(key string) {
	if key == "_index" {
		st.CreateIndex(key)
	}
}

// CreateIndex implementes creational of teh new index
func (st *Store) CreateIndex(index string) {
	if index == "" {
		log.Info(fmt.Sprintf("New index %s can't be created", index))
		return
	}
	st.index.CreateIndex(index)
}

func (st *Store) CheckExistDB(value string) bool {
	_, ok := st.dbs[value]
	return ok
}

func (st *Store) CreateDB(dbname string) {
	st.lock.Lock()
	defer st.lock.Unlock()
	_, ok := st.dbs[dbname]
	st.dbs[dbname] = CreateNewDB(dbname)
	if !ok {
		st.stat.Dbnum += 1
	}
}

// Get provides getting of the value by the key
func (st *Store) Get(value []byte) []byte {
	return st.get(value, "")
}

//if dbname is not equal "", get data from db with name dbname
func (st *Store) get(key []byte, dbname string) []byte {
	if st == nil {
		return nil
	}
	store := st.store
	if store == nil {
		return nil
	}
	if dbname != "" {
		//check db availability
		dbdata, ok := st.dbs[dbname]
		if !ok {
			log.Fatal(fmt.Sprintf("db with name %s is not found", dbname))
		} else if !dbdata.isactive {
			log.Fatal(fmt.Sprintf("db with name %s is not active", dbname))
		}
	}

	st.lock.RLock()
	defer func() {
		st.lock.RUnlock()
		st.stat.NumReads++
	}()
	result, err := store.Get(key)
	if err != nil {
		st.stat.NumFailedReads++
		return nil
	}

	if st.compression {
		result = decompress(result.([]byte))
	}
	fmt.Println(string(result.([]byte)))

	return result.([]byte)

}

//Get many keys from list
func (st *Store) GetMany(keys [][]byte) interface{} {
	result := make([]interface{}, len(keys))
	st.lock.Lock()
	defer st.lock.Unlock()
	if len(keys) > 0 {
		for i := 0; i < len(keys); i++ {
			result = append(result, st.Get(keys[i]))
		}
		st.historyevent.AddEvent("default", "GetMany")
		return result
	}
	return nil
}

//check and split keys on system and not
func (st *Store) beforeSetKeys(items KVITEM) *ReadyToSet {
	//Vefore set, check split system keys with prefix _ and simple keys
	return NewReadyToSet(items)
}

// Set provides setting of key-value pair
func (st *Store) Set(key, value []byte) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	if err := st.beforeSet(key, value); err != nil {
		return err
	}
	if st.compression {
		value = compress(value)
	}
	return st.store.Put(key, value, ds.ItemOptions{})
}

func (st *Store) beforeSet(key, value []byte) error {
	if len(key) > MaxKeySize {
		return errMaxKeySize
	}
	if len(value) > MaxValueSize {
		return errMaxValueSize
	}
	return nil
}

//Before set data to the lightstore. Check if in current request
//exists system keys with prefix _
func getItemOptions(items map[string]string) ds.ItemOptions {
	itemopt := ds.ItemOptions{}
	for key, value := range items {
		if key == "_index" {
			itemopt.Index = value
		}
		if key == "_immutable" {
			itemopt.Immutable = false
			if value == "true" {
				itemopt.Immutable = true
			}
		}
		if key == "_update" {
			itemopt.Update = false
			if value == "true" {
				itemopt.Update = true
			}
		}
	}
	fmt.Println("ITM: ", itemopt)
	return itemopt
}

//Exist check key in the lightstore
//and return true if key exist and false otherwise
func (st *Store) Exist(key []byte) bool {
	store := st.store
	return store.Exist(key)
}

func (st *Store) set(dbname string, key []byte, value []byte, opt ds.ItemOptions) bool {
	st.lock.Lock()
	defer st.lock.Unlock()
	store := st.store
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
			log.Info(fmt.Sprintf("db with name %s not availability to write, because contains maxumum possible number of data objrct", dbname))
			return false
		}
	}

	go func(s ds.Storage) {
		starttime := time.Now()
		s.Put(key, value, ds.ItemOptions{})
		if dbname != "" {
			dbdata, _ := st.dbs[dbname]
			st.historyevent.AddEvent("deafult", "Set")
			dbdata.datacount += 1
		}
		st.PublishKeyValue(string(key), string(dbname))
		st.stat.NumWrites += 1
		fmt.Println(fmt.Sprintf("Stored in : %s", time.Since(starttime)))
	}(store)

	return true

}

//Remove provides clearning curent key
func (st *Store) Remove(key []byte) {
	st.store.Delete(key)
}

func (st *Store) Find(key []byte) interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	return nil
}

// Stat retruns statistics on store
func (st *Store) Stat() *statistics.Statistics {
	return st.stat
}

func (st *Store) Close() {
	fmt.Println("End working: ", time.Now())
}

func (st *Store) SubscribeKey(item string) {
	rpc.InitClient(nil).Get("Pubsub.Subscribe", &SubscribeData{Title: item}, nil)
	//st.pubsub.Subscribe(&SubscribeData{Title: item})
}

func (st *Store) PublishInfo(key string) {
	st.pubsub.Publish(&PublishData{Title: key, Msg: "newmsg"})
}

func (st *Store) PublishKeyValue(key, value string) {
	st.pubsub.Publish(&PublishData{Title: key, Msg: value})
}

// ISCreated returns true of store was created
func (st *Store) IsCreated() bool {
	return true
}

func (st *Store) makeSnapshot() {
}

//This private method provides checking inner datastructure for storing
func checkDS(name string) (result ds.Storage) {
	result = ds.NewDict()
	/* SkipList datastructure as main store */
	if name == "skiplist" {
		result = ds.NewSkipList()
	}

	/* Simple map as main store */
	if name == "dict" {
		result = ds.NewDict()
	}

	/*B-tree structure as main store */
	if name == "b-tree" {
		result = ds.InitBTree(10)
	}
	return result
}

//Every provides doing some operation every n seconds/minutes
//Data for this doing, reading from config
//For example: Doing snapshots every 1 minute
func (store *Store) Every(funcs []func()) {
	go func() {
		for {
			for _, f := range funcs {
				go f()
			}
			time.Sleep(time.Duration(store.config.Every.Seconds) * time.Second)
		}
	}()
}

func InitStore(settings Settings) *Store {
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	mutex := &sync.RWMutex{}
	store := new(Store)
	starttime := time.Now().UTC()
	store.items = 0
	store.store = checkDS(settings.Innerdata)
	store.keys = []string{}
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(statistics.Statistics)
	store.stat.Start = starttime
	store.index = NewIndexing()
	store.config = LoadConfigData("")
	store.pubsub = PubsubInit()
	rpc.RegisterRPCFunction(store.pubsub)
	store.rpcdata = rpc.Init("")
	store.rpcdata.Run()
	return store
}
