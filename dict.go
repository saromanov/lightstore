package lightstore

type Dict struct {
	Value map[string]*DictItem
	stat ItemStatistics
}

type DictItem struct {
	value interface{}
	stat *ItemStatistics
}

func (ditem *DictItem) setValue(value interface{}) {
	ditem.value = value
	ditem.stat = InitItemStatistics()
}

func (ditem *DictItem) getValue() interface{} {
	ditem.stat.num_reads += 1
	return ditem.value
}

func NewDict() *Dict {
	d := new(Dict)
	d.Value = make(map[string]*DictItem)
	return d
}

func (d *Dict) Set(key string, value interface{}) {
	d.Value[key].setValue(value)
}

func (d *Dict) Get(key string) (interface{}, bool) {
	value, ok := d.Value[key]
	return value.getValue(), ok
}

func (d *Dict) Remove(key string) {
	_, ok := d.Value[key]
	if ok {
		delete(d.Value, key)
	}
}
