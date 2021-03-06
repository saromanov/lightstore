package store

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	// TrivialMode defines simple key-value store
	TrivialMode = "Trivial"
	// ServerMode mode creates server for key-value
	ServerMode = "Server"
	// ClusterMode creates cluster for store
	ClusterMode = "Cluster"
)

//This module for loading configuration from config.yaml

// Config defines main config for Lightstore
type Config struct {
	//Address for server
	Address string
	//port for server
	Port uint
	//Directory for storing log data
	Logdir string
	//Directory for storing cache data
	Cachedir string
	//Diretcory for storing db data
	Dbdir string
	//Cluster name
	Cluster string
	//Sync for commits
	Commitsync bool
	//Size for cache
	Cachesize int
	//Sise for blocks
	Blocksize int
	// Storage defines name of the type
	// for inner storage
	Storage string

	Every struct {
		Seconds int
		Actions []string
	}

	//Limit list for the history of events
	Historylimit int

	// Set compression of data
	Compression bool

	// Mode provides setting of the mode of Lightstore
	// Trivial, Server, Cluster
	Mode string

	// If this flag is true then enable Prometheus
	// as monitoring provider
	Monitoring bool

	// MaxKeySize defines maximum key size for store
	MaxKeySize uint

	// MaxValueSize defines maximum valuze size for store
	MaxValueSize uint

	// LoadPath provides definition for load data from disk
	LoadPath string
}

//LoadConfigData provides load configuration or set default params
func LoadConfigData(path string) *Config {
	if path == "" {
		path = getConfigPath()
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return defaultConfig()
	}
	var conf Config
	yamlerr := yaml.Unmarshal(data, &conf)
	if yamlerr != nil {
		panic(yamlerr)
	}

	conf.setMissedValues()
	return &conf
}

func (conf *Config) setMissedValues() {
	if conf.Address == "" {
		conf.Address = "localhost"
	}

	if conf.Port == 0 {
		conf.Port = 8080
	}

	if conf.Logdir == "" {
		conf.Logdir = setDefaultLogPath()
	}

	if conf.Dbdir == "" {
		conf.Dbdir = setDefaultDBData()
	}

	if conf.Historylimit == 0 {
		conf.Historylimit = 1000
	}

	if conf.Storage == "" {
		conf.Storage = "dict"
	}

	if conf.MaxKeySize == 0 {
		conf.MaxKeySize = MaxKeySize
	}

	if conf.MaxValueSize == 0 {
		conf.MaxValueSize = MaxValueSize
	}

	conf.Mode = checkMode(conf.Mode)
}

// checkMode provides setting of Lightstore mode to config
func checkMode(mode string) string {
	if mode != TrivialMode && mode != ServerMode && mode != ClusterMode {
		return TrivialMode
	}
	return mode
}

// defaultConfig creates default attributes for DB
func defaultConfig() *Config {
	conf := new(Config)
	conf.Address = "localhost"
	conf.Port = 8080
	conf.Logdir = setDefaultLogPath()
	conf.Dbdir = setDefaultDBData()
	conf.Cluster = "cluster1"
	conf.Cachesize = 1024
	conf.Storage = "dict"
	conf.Monitoring = false
	conf.MaxKeySize = MaxKeySize
	conf.MaxValueSize = MaxValueSize
	return conf
}

func getConfigPath() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/lightstore/config.yaml", home)
}

//Set log path (/var/log/litghstore)
func setDefaultLogPath() string {
	path := fmt.Sprintf("%s/lightstore", os.Getenv("HOME"))
	filepath := fmt.Sprintf("%s/lightstore.log", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		//Note, need to create lightstore user
		os.Mkdir(path, 0777)
		f, errfile := os.Create(fmt.Sprintf("%s/lightstore.log", path))
		defer f.Close()
		if errfile != nil {
			panic(errfile)
		}

	} else {
		_, errfile := os.Stat(filepath)
		if errfile != nil {
			f, _ := os.Create(fmt.Sprintf("%s/lightstore.log", path))
			defer f.Close()
		}
	}

	return filepath
}

func setDefaultDBData() string {
	home := os.Getenv("HOME")
	path := fmt.Sprintf("%s/lightstore", home)
	filepath := fmt.Sprintf("%s/lightstore.data", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		//Note, need to create lightstore user
		os.Mkdir(path, 0777)
		f, errfile := os.Create(filepath)
		defer f.Close()
		if errfile != nil {
			panic(errfile)
		}
		return filepath

	}
	_, errfile := os.Stat(filepath)
	if errfile != nil {
		f, _ := os.Create(fmt.Sprintf("%s/lightstore.data", path))
		defer f.Close()
	}

	return filepath
}
