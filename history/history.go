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
	lock sync.RWMutex
}

type Event struct{
	//Title of event
	title string
	//Address from event
	addr  string
	//Time where event was happen
	timesdata time.Time
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
	if len(hist.items) == hist.limit {
		hist.removeOutdated(1)
	}

	hist.items = append(hist.items, &Event{
		title: title,
		addr: addr, 
		timesdata: time.Now(),
		})
}

//Get event by id
func (hist*History) Get(idx int)*Event {
	if idx > hist.limit {
		return &Event{}
	}

	return hist.items[idx]
}

func (hist *History) GetAll()[]*Event {
	return hist.items
}

//Remove outdated records from items.
//Records removed from last positions
func (hist*History) removeOutdated(rempos int){
	hist.items = append(hist.items[:0], hist.items[rempos:]...)
}