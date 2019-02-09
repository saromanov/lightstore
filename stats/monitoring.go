package stats

import (
	"sync"
	"time"
)

//Statistics provides basic statistics for all db
type Statistics struct {
	lock *sync.RWMutex
	//Total number of reads
	NumReads int
	//Total number of writes
	NumWrites int
	//Start time
	Start time.Time
	//Number of active db
	Dbnum int
	// NumFailedReads retruns number of failed reads
	NumFailedReads int
}

//This struct provides statistics for each item
type ItemStatistics struct {
	key      string
	NumReads int
	start    time.Time
}

func InitItemStatistics() *ItemStatistics {
	start := time.Now()
	itemstat := new(ItemStatistics)
	itemstat.start = start
	return itemstat
}

//IncrWrites provides increment of total number of writes
func (dbm *Statistics) IncrWrites() {
	dbm.lock.RLock()
	defer dbm.lock.RUnlock()
	dbm.NumWrites++
}

//IncrReads provides increment of total number of reads
func (dbm *Statistics) IncrReads() {
	dbm.lock.RLock()
	defer dbm.lock.RUnlock()
	dbm.NumReads++
}
