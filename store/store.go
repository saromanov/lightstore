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
	// MaxKeySize defines maximim size of the key
	MaxKeySize uint = 512
	// MaxValueSize defines maximim size of the value
	MaxValueSize uint = 32768
)

type Settings struct {
	Innerdata string
}

// Store provides implementation of the main store
type Store struct {
	dbs         map[string]*DB
	store       ds.Storage
	storename   string
	keys        []string
	lock        *sync.RWMutex
	stat        *stats.Statistics
	config      *Config
	compression bool
	fileWatcher *watcher
	writer      *Writer
	indexes     map[string]*index
}

// newStore creates a new instance of lightstore
func newStore(c *Config) (*Store, error) {
	if c == nil {
		c = defaultConfig()
	}
	mutex := &sync.RWMutex{}
	store := new(Store)
	startTime := time.Now().UTC()
	fileWatcher, err := newWatcher(".")
	if err != nil {
		log.Info(fmt.Sprintf("unable to init file watcher: %v", err))
	}
	if c.Monitoring {
		monitoring.Init()
	}
	store.fileWatcher = fileWatcher
	store.store = makeStorage("")
	store.keys = []string{}
	store.compression = c.Compression
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(stats.Statistics)
	store.stat.Start = startTime
	store.indexes = make(map[string]*index)
	c.setMissedValues()
	store.config = c
	if c.LoadPath != "" {
		errLoad := loadData(store, c.LoadPath)
		if errLoad != nil {
			return nil, fmt.Errorf("unable to load data: %v", errLoad)
		}
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
			if commandSet {
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
func (st *Store) CreateIndex(index string) error {
	if index == "" {
		return errNoIndexName
	}
	_, ok := st.indexes[index]
	if ok {
		return fmt.Errorf("index with name %s already exist", index)
	}
	return nil
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

// ISCreated returns true of store was created
func (st *Store) IsCreated() bool {
	return true
}

func (st *Store) makeSnapshot() {
}

// makeStorage provides creating of the engine for storage
func makeStorage(name string) ds.Storage {
	switch name {
	case "skiplist":
		return ds.NewSkipList()
	case "dict":
		return ds.NewDict()
	case "b-tree":
		return ds.InitBTree(10)
	}
	return ds.NewDict()
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
	store.store = makeStorage(settings.Innerdata)
	store.keys = []string{}
	store.dbs = make(map[string]*DB)
	store.lock = mutex
	store.stat = new(stats.Statistics)
	store.stat.Start = startTime
	store.config = LoadConfigData("")
	return store
}
