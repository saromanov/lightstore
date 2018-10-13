package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/ryszard/goskiplist/skiplist"
	ds "github.com/saromanov/lightstore/datastructures"
	"github.com/saromanov/lightstore/history"
	log "github.com/saromanov/lightstore/logging"
	"github.com/saromanov/lightstore/rpc"
	"github.com/saromanov/lightstore/scan"
	"github.com/saromanov/lightstore/snapshot"
	"github.com/saromanov/lightstore/statistics"
	"github.com/saromanov/lightstore/utils"
)

//Basic implmentation of key-value store(without assotiation with any db name)

const param = 0

type Settings struct {
	Innerdata string
}

// Store provides implementation of the main store
type Store struct {
	items         int
	dbs           map[string]*DB
	mainstore     interface{}
	mainstorename string
	keys          []string
	lock          *sync.RWMutex
	stat          *statistics.Statistics
	index         *Indexing
	config        *Config
	pubsub        *Pubsub
	//Event history
	historyevent *history.History
	rpcdata      *rpc.RPCData
}

// Open creates a new instance of lightstore
func Open(c *Config) *Store {
	mutex := &sync.RWMutex{}
	store := new(Store)
	starttime := time.Now().UTC()
	store.items = 0
	store.mainstore = checkDS("")
	store.keys = []string{}
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(statistics.Statistics)
	store.stat.Start = starttime
	store.index = NewIndexing()
	store.config = LoadConfigData("")
	store.ConstructFromConfig()
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
	_, ok := st.dbs[dbname]
	st.dbs[dbname] = CreateNewDB(dbname)
	if !ok {
		st.stat.Dbnum += 1
	}
	st.lock.Unlock()
}

func (st *Store) Get(value []byte) interface{} {
	return st.get(value, "")
}

func (st *Store) GetFromDB(dbname string, value []byte) interface{} {
	return st.get(value, dbname)
}

//if dbname is not equal "", get data from db with name dbname
func (st *Store) get(value []byte, dbname string) interface{} {
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

	st.lock.RLock()
	defer st.lock.RUnlock()
	switch mainstore.(type) {
	case *ds.Dict:
		result, ok := mainstore.(*ds.Dict).Get(value)
		if ok {
			st.stat.NumReads += 1
			st.historyevent.AddEvent("default", "Get")
			return result
		}
	case *skiplist.SkipList:
		result, ok := mainstore.(*skiplist.SkipList).Get(value)
		if ok {
			st.stat.NumReads += 1
			return result
		}
	}

	return nil

}

func (st *Store) AppendData(kvitem KVITEM) {
	for key, value := range kvitem {
		exist := st.Exist(key)
		if exist {
			data := []interface{}{value}
			data = append(data, value)
			st.set("", key, data, ds.ItemOptions{})
		} else {
			items := st.Get(key)
			switch items.(type) {
			case []interface{}:
				items = append(items.([]interface{}), value)
				st.historyevent.AddEvent("default", "Append")
				st.set("", key, items, ds.ItemOptions{})
			default:
				data := []interface{}{st.Get(key)}
				data = append(data, value)
				st.set("", key, data, ds.ItemOptions{})
			}
		}
	}
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
		st.historyevent.AddEvent("default", "GetMany")
		return result
	}
	return nil
}

//check and split keys on system and not
func (st *Store) beforeSet(items KVITEM) *ReadyToSet {
	//Vefore set, check split system keys with prefix _ and simple keys
	return NewReadyToSet(items)
}

func (st *Store) Set(items map[string]string) bool {
	before := st.beforeSet(items)
	if before.ready {
		opt := getItemOptions(before.syskeys)
		for key, value := range before.kvitems {
			st.set("", key, value, opt)
		}
		return true
	} else {
		return false
	}
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
func (st *Store) Exist(key string) bool {
	mainstore := st.mainstore
	switch mainstore.(type) {
	case *ds.Dict:
		return mainstore.(*ds.Dict).Exist(key)
	}
	return false
}

func (st *Store) SetinDB(dbname string, key string, value interface{}) bool {
	return st.set(dbname, key, value, ds.ItemOptions{})
}

func (st *Store) set(dbname string, key string, value interface{}, opt ds.ItemOptions) bool {
	st.lock.Lock()
	defer st.lock.Unlock()
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
			log.Info(fmt.Sprintf("db with name %s not availability to write, because contains maxumum possible number of data objrct", dbname))
			return false
		} else {
			mainstore = dbdata.mainstore
		}
	}

	go func() {
		starttime := time.Now()
		switch mainstore.(type) {
		case *ds.Dict:
			mainstore.(*ds.Dict).Set(key, value, opt)
		case *skiplist.SkipList:
			mainstore.(*skiplist.SkipList).Set(key, value)
		}

		if dbname != "" {
			dbdata, _ := st.dbs[dbname]
			st.historyevent.AddEvent("deafult", "Set")
			dbdata.datacount += 1
		}
		st.PublishKeyValue(key, dbname)
		st.stat.NumWrites += 1
		fmt.Println(fmt.Sprintf("Stored in : %s", time.Since(starttime)))
	}()

	return true

}

//Remove provides clearning curent key
func (st *Store) Remove(key string) {
	switch st.mainstore.(type) {
	case *ds.Dict:
		st.mainstore.(*ds.Dict).Remove(key)
	case *skiplist.SkipList:
		st.mainstore.(*skiplist.SkipList).Delete(key)
	}
}

func (st *Store) Find(key string) interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	scanner := scan.NewScan(key, st.keys)
	if scanner.Find(key) {
		return st.get(key, "")
	}
	return nil
}

func (st *Store) RepairData(key string) *ds.RepairItem {
	fmt.Println(fmt.Sprintf("Try to repair key %s", key))
	item, err := st.mainstore.(*ds.Dict).GetFromRepair(key)
	if err != nil {
		log.Fatal(err.Error())
	}
	return item
}

//Return statistics of usage
func (st *Store) Stat() *statistics.Statistics {
	return st.stat
}

func (st *Store) CloseLightStore() {
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
	//This is only for testing
	snap := snapshot.NewSnapshotObject("")
	snap.Write(&snapshot.SnapshotObject{Crc32: "123", Data: "foobar", Dir: "/"})
}

//This private method provides checking inner datastructure for storing
func checkDS(name string) (result interface{}) {
	result = ds.NewDict()
	/* SkipList datastructure as main store */
	if name == "skiplist" {
		result = skiplist.NewStringMap()
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

//ConstructFromConfig provides creational lightstore params from config
func (store *Store) ConstructFromConfig() {
	if store.config == nil {
		return
	}

	every := store.config.Every
	if len(every.Actions) > 0 {
		store.Every(utils.ActionsNamesToFuncs(every.Actions))
	}

	store.historyevent = history.NewHistory(5)
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
	store.mainstore = checkDS(settings.Innerdata)
	store.keys = []string{}
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(statistics.Statistics)
	store.stat.Start = starttime
	store.index = NewIndexing()
	store.config = LoadConfigData("")
	store.ConstructFromConfig()
	store.pubsub = PubsubInit()
	rpc.RegisterRPCFunction(store.pubsub)
	store.rpcdata = rpc.Init("")
	store.rpcdata.Run()
	return store
}
