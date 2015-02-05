package lightstore

import 
(
	"fmt"
	"sync"
	"github.com/ryszard/goskiplist/skiplist"
	"time"
	//"runtime"
	//"errors"
)


//Basic implmentation of key-value store

const 
(
	param = 0
)

type Settings struct {
	Innerdata string
}

type Dict struct {
	Value map[string] interface{}
}

func Newdict()(*Dict){
	d := new(Dict)
	d.Value = make(map[string] interface{})
	return new (Dict)
}


type Store struct{
	items int
	mainstore interface{}
	mainstorename string
	lock *sync.Mutex
}

type User struct {
    score float64
    id    string
}

type Value struct {
	id float64
	key string
	value interface {}
}
func (st*Store) Get(value string)interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()

	switch st.mainstore.(type){
	case Dict:
		fmt.Println("not implement yet")
	case BMtree:
		fmt.Println("This is bmtree")
	default:
		result, ok := st.mainstore.(*skiplist.SkipList).Get(value)
		if ok {
			return result
		}
	}
	return nil
}

func (st*Store) Set(key string, value interface{}){
	switch st.mainstore.(type){
	case Dict:
		fmt.Println("A")
	default:
		st.mainstore.(*skiplist.SkipList).Set(key, value)
	}
}

func (st*Store) Remove(key string){
	switch st.mainstore.(type){
	case Dict:
		fmt.Println("not implemented yet")
	/*default:
		st.mainstore.(*skiplist.SkipList).Delete(key)*/
	}
}

func (st*Store) CloseLightStore(){
	fmt.Println("End working: ", time.Now())
}

func InitLightStore(settings Settings)(*Store){
	/*
		Type store can be skiplist or b-tree
	*/
	mutex := &sync.Mutex{}
	store := new(Store)
	fmt.Println("Start working: ", time.Now())
	store.items = 0;
	if settings.Innerdata =="skiplist"{
		store.mainstore = skiplist.NewStringMap()
	}

	if settings.Innerdata =="dict"{
		store.mainstore = Newdict()
	}

	if settings.Innerdata == "b-tree" {
		store.mainstore = InitBMTree()
	}
	store.lock = mutex
	return store
}
