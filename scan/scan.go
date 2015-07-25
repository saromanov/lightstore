package scan

// This module provides iteration over all keys in the single collection

import
(
	"sort"
)

//Scan provides iteration over all keys
type Scan struct {
	match string
	keys []string
}

func NewScan(match string, keys []string)*Scan {
	sc := new(Scan)
	sc.match = match
	sc.keys = keys
	return sc
}

func (sc *Scan) Process(){
	sort.Strings(sc.keys)
}

func (sc *Scan) Find(key string) bool {
	for _, item := range sc.keys{
		if item == key {
			return true
		}
	}

	return false
}

func (sc *Scan) Close(){

}

func (sc *Scan) Stop(){
	
}
