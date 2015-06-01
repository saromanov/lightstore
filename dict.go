package lightstore

type Dict struct {
	Value map[string]interface{}
}

func NewDict() *Dict {
	d := new(Dict)
	d.Value = make(map[string]interface{})
	return d
}

func (d *Dict) Set(key string, value interface{}) {

	d.Value[key] = value
}

func (d *Dict) Get(key string) (interface{}, bool) {
	value, ok := d.Value[key]
	return value, ok
}

func (d *Dict) Remove(key string) {
	_, ok := d.Value[key]
	if ok {
		delete(d.Value, key)
	}
}
