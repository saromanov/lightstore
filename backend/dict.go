package backend

import (
	"errors"

	"github.com/saromanov/golib/hashmap"
	"github.com/saromanov/lightstore/stats"
)

type Dict struct {
	engine *hashmap.HashMap
	stat   stats.ItemStatistics
	repair *Repair
}

func NewDict() *Dict {
	d := new(Dict)
	d.engine = hashmap.New()
	d.repair = NewRepair()
	return d
}

func (d *Dict) Put(key []byte, value interface{}, op ItemOptions) error {
	if op.Immutable {
		return nil
	}
	d.engine.Set(key, value)
	return nil
}

// Get provides getting of value by the key
func (d *Dict) Get(key []byte) (interface{}, error) {
	value := d.engine.Get(key)
	if value == nil {
		return nil, errors.New("unable to find element")
	}
	return value, nil
}

func (d *Dict) GetFromRepair(key string) (*RepairItem, error) {
	return d.repair.GetFromRepair(key, "")
}

// Exist provides implementation for checking of key is exist
func (d *Dict) Exist(key []byte) bool {
	exist := d.engine.Get(key)
	if exist == nil {
		return false
	}
	return true
}

// First returns first key from collection
func (d *Dict) First() interface{} {
	return d.engine.FirstKey()
}

// Remove provides removing of the record
func (d *Dict) Delete(key []byte) error {
	d.engine.Remove(key)
	return nil
}

// Next provides iteration over collection
func (d *Dict) Next(i int) interface{} {
	return d.engine.GetKeyByIndex(i)
}
