package lightstore

import
(
	"net/http"
	/*"net/url"
	"time"*/
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"sync"

)
var lock = sync.RWMutex{}
var store = new(Store)

type Item struct{
	Key string
	Value interface{}
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

//Append {Key:value} to dict
func StoreData(w rest.ResponseWriter, r *rest.Request){
	item := Item{}
	err := r.DecodeJsonPayload(&item)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lock.RLock()
	if item.Key == "" {
		rest.Error(w, "Key code required", 400)
		return
	}

	if item.Value == nil {
		rest.Error(w, "Value code required", 400)
		return
	}
	store.Set(item.Key, item.Value)
	lock.RUnlock()

	w.WriteJson("Element was append")
}

//Get key from store and immediately remove
func GetbyKeyAndRemove(w rest.ResponseWriter, r *rest.Request){
	key := r.PathParam("key")
	lock.RLock()
	value := store.Get(key)
	if value == nil {
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


func InitLightStore(typestore string, addr string){
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/get/:key", GetbyKey},
		&rest.Route{"POST", "/set",StoreData},
		&rest.Route{"DELETE", "/remove/:key", DeleteData},
		//Get and delete
		&rest.Route{"POST", "/gad", GetbyKeyAndRemove},
	)

	if err != nil {
		log.Fatal(err)
	}
	store = InitStore(Settings{Innerdata: typestore})
	api.SetApp(router)
	http.ListenAndServe(addr, api.MakeHandler())
}
