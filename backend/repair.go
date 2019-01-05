package backend

import (
	"errors"
	"time"
)

//Repair needs for getting removed data
type Repair struct {
	limit int
	data  []*RepairItem
}

type RepairItem struct {
	Key      string
	Value    interface{}
	Checksum string
	Date     time.Time
}

func NewRepair() *Repair {
	rep := new(Repair)
	rep.data = []*RepairItem{}
	return rep
}

func (rep *Repair) AddToRepair(key string, value interface{}) {
	newitem := &RepairItem{
		Key:      key,
		Value:    value,
		Checksum: Checksum(key),
		Date:     time.Now()}
	rep.data = append(rep.data, newitem)
}

func (rep *Repair) GetFromRepair(key, value string) (*RepairItem, error) {
	for _, item := range rep.data {
		if item.Key == key || item.Value == value {
			return item, nil
		}
	}

	return nil, errors.New("Repair: Element not found")
}
