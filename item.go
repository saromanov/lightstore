package lightstore
import
(
	"time"
)


type Item struct {
	id int64
	checksum string
	value interface{}
	stat *ItemStatistics
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

func NewItem(value interface{})*Item {
	item := new(Item)
	item.value = value
	item.weights = 0
	item.priority = 0
	item.immutable = false
	item.numpastitems = 10
	item.pastitems = []*PastItem{}
	item.checksum = Checksum(value.(string))
	item.writetime = time.Now()
	return item
}

func (itm *Item) UpdateItem(value interface{}){
	if !itm.immutable {
		itm.setToPast()
		itm.value = value
	}
}

func (itm *Item) setToPast(){
	if len(itm.pastitems) < itm.numpastitems {
		newpast := new(PastItem)
		newpast.vernum = len(itm.pastitems)+1
		newpast.updatetime = time.Now()
		newpast.pastitem = itm
		itm.pastitems = append(itm.pastitems, newpast)
	}
}