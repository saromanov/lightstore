package lightstore
import
(
	"time"
	"errors"
)


//Repair needs for getting removed data
type Repair struct {
	limit  int
	data   []*RepairItem
}

type RepairItem struct {
	key string
	value string
	checksum string
	date time.Time
}

func NewRepair()*Repair {
	rep := new(Repair)
	rep.data = []*RepairItem{}
	return rep
}

func (rep*Repair) AddToRepair(key, value string){
	newitem := &RepairItem{
		key:key, 
		value:value, 
		checksum: Checksum(key + value), 
		date:time.Now()}
	rep.data = append(rep.data, newitem)
}

func (rep *Repair) GetFromRepair(key, value string)(*RepairItem, error) {
	for _, item := range rep.data {
		if item.key == key || item.value == value {
			return item, nil
		}
	}

	return nil, errors.New("Repair: Element not found")
}