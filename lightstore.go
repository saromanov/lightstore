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
	stat *Statistics
}

type Statistics struct {
	//Total number of reads
	num_reads int
	//Total number of writes
	num_writes int
}

func (st*Store) Get(value string)interface{} {
	st.lock.Lock()
	defer st.lock.Unlock()
	switch st.mainstore.(type){
	case *Dict:
		result, ok := st.mainstore.(*Dict).Get(value)
		if ok {
			st.stat.num_reads += 1
			return result;
		}
	case *BMtree:
		fmt.Println("Not implemented yet")
	case *skiplist.SkipList:
		result, ok := st.mainstore.(*skiplist.SkipList).Get(value)
		if ok {
			st.stat.num_reads += 1
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
	st.stat.num_writes += 1
}

func (st*Store) Remove(key string){
	switch st.mainstore.(type){
	case Dict:
		st.mainstore.(*Dict).Remove(key)
	default:
		st.mainstore.(*skiplist.SkipList).Delete(key)
	}
}

//Return statistics of usage
func (st*Store) Stat()(int,int){
	return st.stat.num_writes, st.stat.num_reads
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
	store.stat = new(Statistics)
	return store
}
