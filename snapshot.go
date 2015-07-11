package lightstore
import
(
	"encoding/json"
	"io/ioutil"
	"path"
	"fmt"
	"strings"
)

//Basic snapshot for all data in ligtstore

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

	errwrite := ioutil.WriteFile(path.Join(so.Dir, "snap1.lssnapshot"), b, 0777)
	if errwrite != nil {
		panic(err)
	}

}


//Read provides reading snapshot and store data to lightstore
//if snapshotname is ""(empty), load more recently snapshot
func (so *SnapshotObject) Read(snapshotname string) {
	data, err := ioutil.ReadFile(path.Join(so.Dir, snapshotname))
	if err != nil {
		panic(fmt.Sprintf("Can't find snapshot with the name %s", snapshotname))
	}

	snapshots := checkAvailableSnapshots(so.Dir)
}

//Return snapshotnames
func checkAvailableSnapshots(dir string)[]string{
	files, err := ioutil.ReadDir(dir)
	result := []string{}
	if err != nil {
		panic(err)
	}

	for _, pathobj := range files {
		if strings.HasSuffix(pathobj.Name(), "lssnapshot") {
			result = append(result, path.Join(dir, pathobj.Name()))
		}
	}
	return result
}
