package lightstore

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	//Address for server
	address string
	//port for server
	port uint
	//Directory for storing log data
	logdir string
	//Directory for storing cache data
	cachedir string
	//Diretcory for storing db data
	dbdir string
	//Cluster name
	cluster string
	//Sync for commits
	commitsync bool
	//Size for cache
	cachesize int
}


//LoadConfigData provides load configuration or set default params
func LoadConfigData() *Config {
	data, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		fmt.Println("This is")
		setDefaultParams()
	}
	conf := Config{}
	yamlerr := yaml.Unmarshal([]byte(data), &conf)
	if yamlerr != nil {
		panic(yamlerr)
	}

	return setDefaultParams()
}

//In the case if config file is not exist or not full,
// set for each param default value
func setDefaultParams() *Config {
	conf := new(Config)
	conf.address = "localhost"
	conf.port = 8080
	conf.logdir = setDefaultLogPath()
	conf.dbdir = setDefaultDBData()
	conf.cluster = "cluster1"
	conf.cachesize = 1024
	return conf
}

func getConfigPath() string {
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/lightstore/config.yaml", home)
}

//Set log path (/var/log/litghstore)
func setDefaultLogPath() string{
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

func setDefaultDBData() string{
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

	} else {
		_, errfile := os.Stat(filepath)
		if errfile != nil {
			f, _ := os.Create(fmt.Sprintf("%s/lightstore.data", path))
			defer f.Close()
		}
	}

	return filepath
}