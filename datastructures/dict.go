package datastructures

import (
	"github.com/saromanov/golib/hashmap"
	"github.com/saromanov/lightstore/statistics"
)

type Dict struct {
	Value  *hashmap.HashMap
	stat   statistics.ItemStatistics
	repair *Repair
}

func NewDict() *Dict {
	d := new(Dict)
	d.Value = hashmap.New()
	d.repair = NewRepair()
	return d
}

func (d *Dict) Set(key []byte, value interface{}, op ItemOptions) {
	if op.Immutable {
		return
	}
	d.Value.Set(key, value)
}

// Get provides getting of value by the key
// In the case if key is not found, its return nil and false
func (d *Dict) Get(key []byte) (interface{}, bool) {
	value := d.Value.Get(key)
	if value == nil {
		return nil, false
	}
	return value, true
}

func (d *Dict) GetFromRepair(key string) (*RepairItem, error) {
	return d.repair.GetFromRepair(key, "")
}

// Exist provides implementation for checking of key is exist
func (d *Dict) Exist(key []byte) bool {
	exist := d.Value.Get(key)
	if exist == nil {
		return false
	}
	return true
}

// Remove provides removing of the record
func (d *Dict) Delete(key []byte) error {
	d.Value.Remove(key)
	return nil
}
