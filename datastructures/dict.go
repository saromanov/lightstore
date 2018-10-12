package datastructures

import "github.com/saromanov/lightstore/statistics"

type Dict struct {
	Value  map[string]*Item
	stat   statistics.ItemStatistics
	repair *Repair
}

func NewDict() *Dict {
	d := new(Dict)
	d.Value = make(map[string]*Item)
	d.repair = NewRepair()
	return d
}

func (d *Dict) Set(key string, value interface{}, op ItemOptions) {
	_, ok := d.Value[key]
	if !ok {
		d.Value[key] = NewItem(value)
	}
	if ok && op.Update && !op.Immutable {
		d.Value[key].UpdateItem(value)
	}
	if ok && !op.Immutable {
		d.Value[key] = NewItem(value)
	} else {

	}
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
