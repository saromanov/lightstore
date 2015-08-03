package lightstore
import
(
	"time"
	"errors"
	"fmt"
)


//Repair needs for getting removed data
type Repair struct {
	limit  int
	data   []*RepairItem
}

type RepairItem struct {
	key string
	value interface{}
	checksum string
	date time.Time
}

func NewRepair()*Repair {
	rep := new(Repair)
	rep.data = []*RepairItem{}
	return rep
}

func (rep*Repair) AddToRepair(key string, value interface{}){
	newitem := &RepairItem{
		key:key, 
		value:value, 
		checksum: Checksum(key), 
		date:time.Now()}
	rep.data = append(rep.data, newitem)
}

func (rep *Repair) GetFromRepair(key, value string)(*RepairItem, error) {
	for _, item := range rep.data {
		fmt.Println(item.key)
		if item.key == key || item.value == value {
			return item, nil
		}
	}

	return nil, errors.New("Repair: Element not found")
}