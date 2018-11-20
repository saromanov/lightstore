package store

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//This module for loading configuration from config.yaml

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
}

//LoadConfigData provides load configuration or set default params
func LoadConfigData(path string) *Config {
	if path == "" {
		path = getConfigPath()
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return setDefaultParams()
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
}

//In the case if config file is not exist or not full,
// set for each param default value
func setDefaultParams() *Config {
	conf := new(Config)
	conf.Address = "localhost"
	conf.Port = 8080
	conf.Logdir = setDefaultLogPath()
	conf.Dbdir = setDefaultDBData()
	conf.Cluster = "cluster1"
	conf.Cachesize = 1024
	conf.Storage = "dict"
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
