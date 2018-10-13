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

func (d *Dict) Get(key string) (interface{}, bool) {
	_, ok := d.Value[key]
	if !ok {
		return nil, ok
	} else {
		d.Value[key].stat.NumReads += 1
		return d.Value[key].value, ok
	}
}

func (d *Dict) GetFromRepair(key string) (*RepairItem, error) {
	return d.repair.GetFromRepair(key, "")
}

func (d *Dict) Exist(key string) bool {
	_, ok := d.Value[key]
	return ok
}

func (d *Dict) Remove(key string) {
	item, ok := d.Value[key]
	if ok {
		d.repair.AddToRepair(key, item.value)
		delete(d.Value, key)
	}
}
