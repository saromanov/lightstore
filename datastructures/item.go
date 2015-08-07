package datastructures
import
(
	"time"
	"../statistics"
)


type Item struct {
	id int64
	checksum string
	value interface{}
	stat *statistics.ItemStatistics
	weights int
	priority int
	immutable bool
	writetime time.Time
	//Maximum Number of past values
	numpastitems int
	pastitems []*PastItem
}

type PastItem struct {
	//Number of version
	vernum int
	//Updated time
	updatetime time.Time
	//Past item
	pastitem *Item
}

//ItemOptions provides additional options for storing values
type ItemOptions struct {
	immutable bool
	update bool
	index string
}

//NewItem provides creates new item before store in memory or write on disk
func NewItem(value interface{})*Item {
	item := new(Item)
	item.value = value
	item.weights = 0
	item.priority = 0
	item.immutable = false
	item.numpastitems = 10
	item.stat = statistics.InitItemStatistics()
	item.pastitems = []*PastItem{}
	item.checksum = Checksum(value.(string))
	item.writetime = time.Now()
	return item
}

//UpdateItem provides set new item for exist in the case, 
//if this item is not immutable
func (itm *Item) UpdateItem(value interface{}){
	if !itm.immutable {
		itm.setToPast()
		itm.value = value
	}
}


//set to past version of the item
func (itm *Item) setToPast(){
	if len(itm.pastitems) < itm.numpastitems {
		newpast := new(PastItem)
		newpast.vernum = len(itm.pastitems)+1
		newpast.updatetime = time.Now()
		newpast.pastitem = itm
		itm.pastitems = append(itm.pastitems, newpast)
	}
}