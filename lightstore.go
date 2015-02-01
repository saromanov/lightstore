package main

import 
(
	"fmt"
	"sync"
	"github.com/gansidui/skiplist"
	"time"
	//"runtime"
	//"errors"
)
const 
(
	param = 0
)

type Settings struct {
	innerdata string
}

type Dict struct {
	value map[string] interface{}
}

func Newdict()(*Dict){
	return new (Dict)
}

func Insert(key string, data interface{}){

}

func Get(key string){
	
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


func (st*Store) Get(value string) {
	st.lock.Lock()
	defer st.lock.Unlock()

	switch st.mainstore.(type){
	case Dict:
		fmt.Println("not implement yet")
	default:
		fmt.Println(st.mainstore.(*skiplist.SkipList).GetElementByRank(1).Value)
	}
}

func (st*Store) Set(key string){
	switch st.mainstore.(type){
	case Dict:
		fmt.Println("A")
	default:
		ust := &User{0.6, key}
		st.mainstore.(*skiplist.SkipList).Insert(ust)
	}
}

func (st*Store) Remove(key string){
	switch st.mainstore.(type){
	case Dict:
		fmt.Println("not implemented yet")
	default:
		st.mainstore.(*skiplist.SkipList).Delete(key)
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
	if settings.innerdata =="skiplist"{
		store.mainstore = skiplist.New()
	}

	if settings.innerdata =="dict"{
		store.mainstore = Newdict()
	}

	if settings.innerdata == "b-tree" {
		
	}
	store.lock = mutex
	return store
}


func (u *User) Less(other interface{}) bool {
    if u.score > other.(*User).score {
        return true
    }
    if u.score == other.(*User).score && len(u.id) > len(other.(*User).id) {
        return true
    }
    return false
}

func test_sync(){
	mutex := &sync.Mutex{}
	go func(){
		total := 0
		for i := 0; i < 100; i++ {
			mutex.Lock()
			total += 1
			mutex.Unlock()
		}
	}()
}

func main() {
	/*fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.NumGoroutine())*/
	/*us := make([]*User, 7)
    us[0] = &User{6.6, "hi"}
    us[1] = &User{3.1, "hi"}
    us[2] = &User{4.5, "hi"}
    us[3] = &User{7.3, "hi"}
	sl := skiplist.New()
	sl.Insert(us[0])
	sl.Insert(us[1])
	sl.Insert(us[2])
	sl.Insert(us[3])*/

	st:= InitLightStore(Settings{innerdata: "skiplist"})
	st.Set("first")
	st.Set("New value")
	st.Get("first")
	st.CloseLightStore()
}