package lightstore

type Dict struct {
	Value map[string]*Item
	stat ItemStatistics
}


func NewDict() *Dict {
	d := new(Dict)
	d.Value = make(map[string]*Item)
	return d
}

func (d *Dict) Set(key string, value interface{}, op ItemOptions) {
	_, ok := d.Value[key]
	if ok && op.update && !op.immutable{
		d.Value[key].UpdateItem(value)
	} else {
		d.Value[key] = NewItem(value)
	}
}

func (d *Dict) Get(key string) (interface{}, bool) {
	_, ok := d.Value[key]
	if !ok {
		return nil, ok
	} else {
		d.Value[key].stat.num_reads += 1
		return d.Value[key].value, ok
	}
}

func (d *Dict) Remove(key string) {
	_, ok := d.Value[key]
	if ok {
		delete(d.Value, key)
	}
}
