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


type Store struct{
	items int
	mainstore interface{}
	mainstorename string
	lock *sync.Mutex
}

func (st*Store) Get(value string)interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	switch st.mainstore.(type){
	case *Dict:
		result, ok := st.mainstore.(*Dict).Get(value)
		if ok {
			return result;
		}
	case *BMtree:
		fmt.Println("Not implemented yet")
	case *skiplist.SkipList:
		result, ok := st.mainstore.(*skiplist.SkipList).Get(value)
		if ok {
			return result
		}
	}
	return nil
}

func (st*Store) Set(key string, value interface{}){
	switch st.mainstore.(type){
	case *Dict:
		st.mainstore.(*Dict).Set(key, value)
	case *skiplist.SkipList:
		st.mainstore.(*skiplist.SkipList).Set(key, value)
	}
}

func (st*Store) Remove(key string){
	switch st.mainstore.(type){
	case Dict:
		st.mainstore.(*Dict).Remove(key)
	default:
		st.mainstore.(*skiplist.SkipList).Delete(key)
	}
}

func (st*Store) CloseLightStore(){
	fmt.Println("End working: ", time.Now())
}

func InitStore(settings Settings)(*Store){
	/*
		Type store can be skiplist or b-tree or simple dict
	*/
	mutex := &sync.Mutex{}
	store := new(Store)
	fmt.Println("Start working: ", time.Now())
	store.items = 0;
	/* SkipList datastructure as main store */
	if settings.Innerdata =="skiplist"{
		store.mainstore = skiplist.NewStringMap()
	}

	/* Simple map as main store */
	if settings.Innerdata =="dict"{
		store.mainstore = NewDict()
	}

	/*B-tree structure as main store */
	if settings.Innerdata == "b-tree" {
		store.mainstore = InitBMTree()
	}
	store.lock = mutex
	return store
}
