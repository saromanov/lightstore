package store

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	ds "github.com/saromanov/lightstore/backend"
	log "github.com/saromanov/lightstore/logging"
	"github.com/saromanov/lightstore/monitoring"
	"github.com/saromanov/lightstore/stats"
)

//Basic implmentation of key-value store(without assotiation with any db name)

const (
	MaxKeySize   uint = 512
	MaxValueSize uint = 32768
)

type Settings struct {
	Innerdata string
}

// Store provides implementation of the main store
type Store struct {
	items       int
	dbs         map[string]*DB
	store       ds.Storage
	storename   string
	keys        []string
	lock        *sync.RWMutex
	stat        *stats.Statistics
	index       *Indexing
	config      *Config
	pubsub      *Pubsub
	compression bool
	fileWatcher *watcher
	writer      *Writer
	indexes     [string]*index
}

// newStore creates a new instance of lightstore
func newStore(c *Config) (*Store, error) {
	if c == nil {
		c = defaultConfig()
	}
	mutex := &sync.RWMutex{}
	store := new(Store)
	startTime := time.Now().UTC()
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
	store.stat = new(stats.Statistics)
	store.stat.Start = startTime
	store.index = NewIndexing()
	store.indexes = make(map[string]index)
	c.setMissedValues()
	store.config = c
	if c.LoadPath != "" {
		err := loadData(store, c.LoadPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load data: %v", err)
		}
	} else {
		c.LoadPath = "lightstore.db"
	}
	store.writer, err = newWriter(c.LoadPath)
	if err != nil {
		return nil, fmt.Errorf("unable to create writer: %v", err)
	}
	return store, nil
}

// loadData provides loading of data from path
func loadData(st *Store, path string) error {
	if path == "" {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var key, value []byte
	var commandSet bool
	var inc int
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.Compare(line, []byte("set;")) == 0 {
			commandSet = true
		} else if bytes.Compare(line, []byte("end;")) == 0 {
			if commandSet == true {
				commandSet = false
				inc = 0
				err := st.store.Put(key, value, ds.ItemOptions{})
				if err != nil {
					return fmt.Errorf("unable to set data: %v", err)
				}
			}
		} else {
			switch inc {
			case 0:
				key = line
				inc++
			case 1:
				value = line

			}
		}
	}
	return nil
}

// CreateIndex implements creational of the new index
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
		st.stat.Dbnum++
	}
}

// Get provides getting of the value by the key
func (st *Store) Get(key []byte) []byte {
	return st.get(key, "")
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
		return result
	}
	return nil
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
	err := st.store.Put(key, value, ds.ItemOptions{})
	if err != nil {
		return err
	}

	return st.writeToLogFile(key, value)
}

func (st *Store) writeToLogFile(key, value []byte) error {
	return st.writer.AddSetCommand(key, value)
}

func (st *Store) beforeSet(key, value []byte) error {
	if uint(len(key)) > MaxKeySize {
		return errMaxKeySize
	}
	if uint(len(value)) > MaxValueSize {
		return errMaxValueSize
	}
	return nil
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
		startTime := time.Now()
		s.Put(key, value, ds.ItemOptions{})
		if dbname != "" {
			dbdata, _ := st.dbs[dbname]
			dbdata.datacount += 1
		}
		st.PublishKeyValue(string(key), string(dbname))
		st.stat.NumWrites += 1
		fmt.Println(fmt.Sprintf("Stored in : %s", time.Since(startTime)))
	}(store)

	return true

}

//Remove provides clearning current key
func (st *Store) Remove(key []byte) {
	st.store.Delete(key)
}

func (st *Store) first() []byte {
	if st == nil {
		return nil
	}
	item := st.store.First()
	if item == nil {
		return nil
	}
	return st.store.First().([]byte)
}

func (st *Store) next(i int) []byte {
	response := st.store.Next(i)
	if response == nil {
		return nil
	}
	return response.([]byte)
}

func (st *Store) Find(key []byte) interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	return nil
}

// Stat retruns statistics on store
func (st *Store) Stat() *stats.Statistics {
	return st.stat
}

func (st *Store) Close() {
	fmt.Println("End working: ", time.Now())
}

func (st *Store) SubscribeKey(item string) {

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
	startTime := time.Now().UTC()
	store.items = 0
	store.store = checkDS(settings.Innerdata)
	store.keys = []string{}
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(stats.Statistics)
	store.stat.Start = startTime
	store.index = NewIndexing()
	store.config = LoadConfigData("")
	store.pubsub = PubsubInit()
	return store
}
