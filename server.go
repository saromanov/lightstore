package main

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

func GetbyKey(w rest.ResponseWriter, r *rest.Request){
	lock.RLock()
	defer lock.RUnlock()
}

func StoreData(w rest.ResponseWriter, r *rest.Request){

}

func InitServer(){
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/get", GetbyKey},
		&rest.Route{"POST", "/set",StoreData},
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.ListenAndServe(":8080", api.MakeHandler())
}
