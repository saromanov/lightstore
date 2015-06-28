package lightstore

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var lock = sync.RWMutex{}
var store = new(Store)
var log = logging.MustGetLogger("lightstore_log")

func checkData(value string) bool {
	if value == "" {
		log.Error("Value code required")
		return false
	}

	return true
}

func GetbyKey(w rest.ResponseWriter, r *rest.Request) {

	key := r.PathParam("key")
	lock.RLock()
	value := store.Get(key)
	if value == nil {
		rest.Error(w, "Value code required", 400)
		return
	}
	defer lock.RUnlock()

	w.WriteJson(value)
}

func GetbyKeyFromDB(w rest.ResponseWriter, r *rest.Request) {
	dbname := r.PathParam("db")
	key := r.PathParam("key")
	lock.RLock()
	value := store.GetFromDB(dbname, key)
	if value == nil {
		rest.Error(w, "Value code required", 400)
		return
	}
	defer lock.RUnlock()

	w.WriteJson(value)
}

func GetbyKeyMany(w rest.ResponseWriter, r *rest.Request) {

}

//Append {Key:value} to dict
func StoreData(w rest.ResponseWriter, r *rest.Request) {
	mapdata := prepareDataToAppend(r.Body)
	lock.RLock()
	go func() {
		log.Info("Store data")
		store.Set(mapdata)
	}()
	defer lock.RUnlock()

	w.WriteJson("Element was append")
}

//This private method contain translation from json to string data
func prepareDataToAppend(r io.ReadCloser) (result map[string]string) {
	var res interface{}
	body, _ := ioutil.ReadAll(r)
	errunm := json.Unmarshal([]byte(body), &res)
	if errunm != nil {
		log.Error(errunm.Error())
		//rest.Error(w, errunm.Error(), http.StatusInternalServerError)
		return
	}
	halfdecoded := res.(map[string]interface{})
	result = make(map[string]string)
	for key, _ := range halfdecoded {
		if checkData(halfdecoded[key].(string)) {
			result[key] = halfdecoded[key].(string)
		}
	}
	return result
}

func StoreDataToDB(w rest.ResponseWriter, r *rest.Request) {
	dbname := r.PathParam("db")
	/*decoder := json.NewDecoder(r.Body)*/
	mapdata := prepareDataToAppend(r.Body)
	lock.RLock()
	go func() {
		log.Info(fmt.Sprintf("Store data in db %s", dbname))
		for key := range mapdata {
			store.SetinDB(dbname, key, mapdata[key])
		}
	}()
	defer lock.RUnlock()

	w.WriteJson(fmt.Sprintf("Element was append in db %s", dbname))

}

func CreateDB(w rest.ResponseWriter, r *rest.Request) {
	log.Info("Try to create new db")
	db := r.PathParam("db")
	store.CreateDB(db)
	w.WriteJson(fmt.Sprintf("db with the name %s was created", db))
}

//Create new index
func CreateIndex(w rest.ResponseWriter, r *rest.Request) {
	indextitle := r.PathParam("index")
	lock.RLock()
	defer lock.RUnlock()
	go func(){
		store.CreateIndex(indextitle)
	}()
	w.WriteJson("Index was created")
}

//Get key from store and immediately remove
func GetbyKeyAndRemove(w rest.ResponseWriter, r *rest.Request) {
	log.Info("Try to GetbyKeyAndRemove")
	key := r.PathParam("key")
	lock.RLock()
	value := store.Get(key)
	if value == nil {
		log.Error("Value code required")
		rest.Error(w, "Value code required", 400)
		return
	}
	store.Remove(key)
	lock.RUnlock()
	w.WriteJson(value)
}

func DeleteData(w rest.ResponseWriter, r *rest.Request) {
	key := r.PathParam("key")
	lock.RLock()
	store.Remove(key)
	lock.RUnlock()
	w.WriteJson("Element was removed")
}

//Append data to store by key
func AppendData(w rest.ResponseWriter, r *rest.Request) {
	/*item := Item{}
	err := r.DecodeJsonPayload(&item)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/

	mapdata := prepareDataToAppend(r.Body)

	lock.RLock()
	store.AppendData(mapdata)
	defer lock.RUnlock()
	w.WriteJson("Data was append to list")
}

//Return statistics of usage
func Show_Statistics(w rest.ResponseWriter, r *rest.Request) {
	lock.Lock()
	log.Info("Try to getting statistics")
	stat := store.Stat()
	w.WriteJson(map[string]int{"Total number of writes": stat.num_writes, "Total number of reads": stat.num_reads,
		"Total number of active db": stat.dbnum})
	lock.Unlock()
}

//Find by key
func Find(w rest.ResponseWriter, r *rest.Request) {

}

//PingPong provides availability of server
func PingPong(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("Pong")
}

func LogConfigure(path string) {
	logfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("Log file as not created")
	}
	logging.SetFormatter(logging.GlogFormatter)
	logbackend := logging.NewLogBackend(logfile, "", 0)
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0), logbackend)
}

func InitLightStore(typestore string, addr string, port uint) {
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	LogConfigure("lightstore.log")
	log.Info("Start to create basic API")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/get/:key", GetbyKey},
		&rest.Route{"GET", "/dbget/:db/:key", GetbyKeyFromDB},
		&rest.Route{"POST", "/set", StoreData},
		&rest.Route{"POST", "/create/:db", CreateDB},
		&rest.Route{"POST", "/createindex/:index", CreateIndex},
		&rest.Route{"POST", "/set/:db", StoreDataToDB},
		&rest.Route{"DELETE", "/remove/:key", DeleteData},
		//Get and delete
		&rest.Route{"POST", "/gad", GetbyKeyAndRemove},
		//Append data to list
		&rest.Route{"POST", "/append", AppendData},
		//Ping the server
		&rest.Route{"GET", "/ping", PingPong},
		//Return short statistics
		&rest.Route{"GET", "/_stat", Show_Statistics},
	)

	if err != nil {
		log.Fatal(err)
	}
	store = InitStore(Settings{Innerdata: typestore})
	api.SetApp(router)
	log.Info("Lightstore is running")
	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), api.MakeHandler())
}
