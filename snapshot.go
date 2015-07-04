package lightstore
import
(
	"encoding/json"
	"io/ioutil"
	"path"
)

//Basic snapshot for the data

type snapshot interface {
	Write(object SnapshotObject)
}

type SnapshotObject struct{
	Crc32 string
	Data  string
	Dir   string
}

func NewSnapshotObject() *SnapshotObject {
	so := new(SnapshotObject)
	so.Dir = "."
	return so
}

func (so* SnapshotObject) Write(object SnapshotObject){
	b, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	errwrite := ioutil.WriteFile(path.Join(so.Dir, "snap1"), b, 0777)
	if errwrite != nil {
		panic(err)
	}

}

func Read(snapshotname string) {

}

