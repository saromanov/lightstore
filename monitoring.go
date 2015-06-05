package lightstore
import
(
	"time"
	"sync"
)

//This module provides monitoring of lightstore


type DBMonitoring struct {
	dbstat *Statistics
	serverstat *ServerStat
	lock *sync.RWMutex
}

//This struct provides basic statistics for all db
type Statistics struct {
	//Total number of reads
	num_reads int
	//Total number of writes
	num_writes int
	//Start time
	start time.Time
	//Number of active db
	dbnum int
}


func InitDBMonitoring() *DBMonitoring {
	start := time.Now()
	return &DBMonitoring {
		&Statistics{0,0, start, 0}, 
		&ServerStat{start, 0},
		&sync.RWMutex{},
	 }
}


//IncrWrite provides increment of total number of writes
func (dbm*DBMonitoring) IncrWrites(){
	dbm.lock.RLock()
	dbm.dbstat.num_writes += 1
	dbm.lock.RUnlock()
}

//IncrRead provides increment of total number of reads
func (dbm*DBMonitoring) IncrReads() {
	dbm.lock.RLock()
	dbm.dbstat.num_reads += 1
	dbm.lock.RUnlock()
}


//This struct provides information about server
type ServerStat struct {
	time_alive time.Time
	numproblems int
}
