package history
import
(
	"time"
	"sync"
)

//This module provides history of events

type History struct {
	items []*Event
	limit int
	count int
	lock sync.RWMutex
}

type Event struct{
	//Title of event
	Title string
	//Address from event
	Addr  string
	//Time where event was happen
	Timesdata time.Time
}

func NewHistory(limit int)* History{
	hist := new(History)
	hist.limit = limit
	hist.lock = sync.RWMutex{}
	hist.items = make([]*Event, limit)
	return hist
}

//AddEvent provides storing new event to list
func (hist*History) AddEvent(addr, title string){
	hist.lock.RLock()
	defer hist.lock.RUnlock()
	if hist.count == hist.limit {
		hist.removeOutdated(1)
		newlog := make([]*Event, hist.limit)
		copy(newlog, hist.items)
		hist.items = newlog
		hist.count = hist.limit-1
	}

	hist.items[hist.count] = &Event{
		Title: title,
		Addr: addr, 
		Timesdata: time.Now(),
		}
	hist.count++
}

//Get event by id
func (hist*History) Get(idx int)*Event {
	hist.lock.RLock()
	defer hist.lock.RUnlock()
	if idx > hist.limit {
		return &Event{}
	}

	return hist.items[idx]
}

func (hist *History) GetAll()[]*Event {
	hist.lock.RLock()
	defer hist.lock.RUnlock()
	return hist.items[0:hist.count]
}

//Remove outdated records from items.
//Records removed from last positions
func (hist*History) removeOutdated(rempos int){
	hist.items = append(hist.items[:0], hist.items[rempos:]...)
}