package lightstore

import
(
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/op/go-logging"
	"sync"
	"os"

)
var lock = sync.RWMutex{}
var store = new(Store)
var log = logging.MustGetLogger("lightstore_log")

type Item struct{
	Key string
	Value interface{}
}

func checkData(w rest.ResponseWriter, item Item) bool {
	if item.Key == "" {
		rest.Error(w, "Key code required", 400)
		return false
	}

	if item.Value == nil {
		rest.Error(w, "Value code required", 400)
		return false
	}

	return true
}


func GetbyKey(w rest.ResponseWriter, r *rest.Request){

	key := r.PathParam("key")
	lock.RLock()
	value := store.Get(key)
	if value == nil {
		rest.Error(w, "Value code required", 400)
		return
	}
	lock.RUnlock()

	w.WriteJson(value)
}

func GetbyKeyMany(w rest.ResponseWriter, r *rest.Request) {
	
}

//Append {Key:value} to dict
func StoreData(w rest.ResponseWriter, r *rest.Request){
	item := Item{}
	err := r.DecodeJsonPayload(&item)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lock.RLock()
	go func(){
		if checkData(w, item){
			log.Info("Store data")
			store.Set(item.Key, item.Value)
		}
	}()
	lock.RUnlock()

	w.WriteJson("Element was append")
}

//Get key from store and immediately remove
func GetbyKeyAndRemove(w rest.ResponseWriter, r *rest.Request){
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

func DeleteData(w rest.ResponseWriter, r *rest.Request){
	key := r.PathParam("key")
	lock.RLock()
	store.Remove(key)
	lock.RUnlock()
	w.WriteJson("Element was removed")
}

//Append data to store by key
func AppendData(w rest.ResponseWriter, r *rest.Request){
	item := Item{}
	err := r.DecodeJsonPayload(&item)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lock.RLock()
	value := store.Get(item.Key)
	if value == nil {
		data := []interface{}{item.Value}
		data = append(data, item.Value)
		store.Set(item.Key, data)
		w.WriteJson("Data was append and create")
	} else {
		items := store.Get(item.Key)
		switch items.(type){
		case []interface{}:
			items = append(items.([]interface{}), item.Value)
			store.Set(item.Key, items)
			w.WriteJson("Data was append to list")
		default:
			data := []interface{}{store.Get(item.Key)}
			data = append(data, item.Value)
			store.Set(item.Key, data)
			w.WriteJson("Data was append and new list was created")
		}
	}

	lock.RUnlock()
}


//Return statistics of usage
func Show_Statistics(w rest.ResponseWriter, r *rest.Request){
	lock.RLock()
	log.Info("Try to getting statistics")
	stat := store.Stat()
	w.WriteJson(map[string]int{"Total number of writes": stat.num_writes, "Total number of reads": stat.num_reads})
	lock.RUnlock()
}


//Find by key
func Find(w rest.ResponseWriter, r *rest.Request) {

}

//PingPong provides availability of server
func PingPong(w rest.ResponseWriter, r *rest.Request){
	w.WriteJson("Pong")
}


func LogConfigure(path string){
	logfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY,0664)
	if err != nil {
		panic("Log file as not created")
	}
	logging.SetFormatter(logging.GlogFormatter)
	logbackend := logging.NewLogBackend(logfile, "",0)
	logging.SetBackend(logging.NewLogBackend(os.Stdout,"",0), logbackend)
}
func InitLightStore(typestore string, addr string){
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	LogConfigure("lightstore.log")
	log.Info("Start to create basic API")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/get/:key", GetbyKey},
		&rest.Route{"POST", "/set",StoreData},
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
	http.ListenAndServe(addr, api.MakeHandler())
}
