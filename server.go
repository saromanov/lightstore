package lightstore

import
(
	"net/http"
	/*"net/url"
	"time"*/
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"sync"
	"fmt"

)
var lock = sync.RWMutex{}
var store = new(Store)

type Item struct{
	key string
	value interface{}
}

func GetbyKey(w rest.ResponseWriter, r *rest.Request){
	lock.RLock()
	defer lock.RUnlock()

	key := r.PathParam("key")

	w.WriteJson(store.Get(key))
}

func StoreData(w rest.ResponseWriter, r *rest.Request){
	item := Item{}
	err := r.DecodeJsonPayload(&item)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(item.key)
	w.WriteJson(&item)
}

func DeleteData(w rest.ResponseWriter, r *rest.Request){

}

func InitServer(typestore string){
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/get", GetbyKey},
		&rest.Route{"POST", "/set",StoreData},
		&rest.Route{"POST", "/remove", DeleteData},
	)

	if err != nil {
		log.Fatal(err)
	}
	store = InitLightStore(Settings{Innerdata: typestore})
	api.SetApp(router)
	http.ListenAndServe(":8080", api.MakeHandler())
}
