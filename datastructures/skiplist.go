package datastructures

import (
	"errors"

	"github.com/ryszard/goskiplist/skiplist"
	"github.com/saromanov/lightstore/statistics"
)

type SkipList struct {
	Value  *skiplist.SkipList
	stat   statistics.ItemStatistics
	repair *Repair
}

func NewSkipList() *SkipList {
	d := new(SkipList)
	d.Value = skiplist.NewStringMap()
	d.repair = NewRepair()
	return d
}

func (d *SkipList) Put(key []byte, value interface{}, op ItemOptions) error {
	if op.Immutable {
		return nil
	}
	d.Value.Set(key, value)
	return nil
}

// Get provides getting of value by the key
func (d *SkipList) Get(key []byte) (interface{}, error) {
	value, ok := d.Value.Get(key)
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
	_, ok := d.Value.Get(key)
	if !ok {
		return false
	}
	return true
}

// Remove provides removing of the record
func (d *SkipList) Delete(key []byte) error {
	_, ok := d.Value.Delete(key)
	if !ok {
		return errors.New("unable to delete element")
	}
	return nil
}
