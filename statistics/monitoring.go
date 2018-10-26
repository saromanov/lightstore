package statistics

import (
	"sync"
	"time"
)

//This module provides monitoring of lightstore

type DBMonitoring struct {
	dbstat     *Statistics
	serverstat *ServerStat
	lock       *sync.RWMutex
}

//This struct provides basic statistics for all db
type Statistics struct {
	//Total number of reads
	NumReads int
	//Total number of writes
	NumWrites int
	//Start time
	Start time.Time
	//Number of active db
	Dbnum int
	//address of host
	Host string
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

func InitDBMonitoring() *DBMonitoring {
	start := time.Now()
	return &DBMonitoring{
		&Statistics{0, 0, start, 0, "localhost:8080", 0},
		&ServerStat{start, 0},
		&sync.RWMutex{},
	}
}

//IncrWrite provides increment of total number of writes
func (dbm *DBMonitoring) IncrWrites() {
	dbm.lock.RLock()
	dbm.dbstat.NumWrites += 1
	dbm.lock.RUnlock()
}

//IncrRead provides increment of total number of reads
func (dbm *DBMonitoring) IncrReads() {
	dbm.lock.RLock()
	dbm.dbstat.NumReads += 1
	dbm.lock.RUnlock()
}

//This struct provides information about server
type ServerStat struct {
	time_alive  time.Time
	numproblems int
}
