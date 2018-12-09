package datastructures

import (
	"errors"

	"github.com/ryszard/goskiplist/skiplist"
	"github.com/saromanov/lightstore/statistics"
)

type SkipList struct {
	engine *skiplist.SkipList
	stat   statistics.ItemStatistics
	repair *Repair
}

func NewSkipList() *SkipList {
	d := new(SkipList)
	d.engine = skiplist.NewStringMap()
	d.repair = NewRepair()
	return d
}

func (d *SkipList) Put(key []byte, value interface{}, op ItemOptions) error {
	if op.Immutable {
		return nil
	}
	d.engine.Set(key, value)
	return nil
}

// Get provides getting of value by the key
func (d *SkipList) Get(key []byte) (interface{}, error) {
	value, ok := d.engine.Get(key)
	if !ok {
		return nil, errors.New("unable to find element")
	}
	return value, nil
}

func (d *SkipList) GetFromRepair(key string) (*RepairItem, error) {
	return d.repair.GetFromRepair(key, "")
}

// Exist provides implementation for checking of key is exist
func (d *SkipList) Exist(key []byte) bool {
	_, ok := d.engine.Get(key)
	if !ok {
		return false
	}
	return true
}

// First returns first element on storage
func (d *SkipList) First() interface{} {
	return d.engine.SeekToFirst().Key()
}

// Remove provides removing of the record
func (d *SkipList) Delete(key []byte) error {
	_, ok := d.engine.Delete(key)
	if !ok {
		return errors.New("unable to delete element")
	}
	return nil
}
