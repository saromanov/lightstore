package snapshot
import
(
	"encoding/json"
	"io/ioutil"
	"path"
	"fmt"
	"strings"
	"os"
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

//NewSnapsjot object provides initialization od new snapshot
func NewSnapshotObject() *SnapshotObject {
	so := new(SnapshotObject)
	so.Dir = "."
	return so
}

//Write provides storing of new snapshot
func (so* SnapshotObject) Write(object* SnapshotObject){
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
	_, err := ioutil.ReadFile(path.Join(so.Dir, snapshotname))
	if err != nil {
		panic(fmt.Sprintf("Can't find snapshot with the name %s", snapshotname))
	}

	snapshots := checkAvailableSnapshots(so.Dir)
	if len(snapshots) == 0{
		log.Info("Can't find available snapshots")
	} else {

	}
}

//Read newest provides reading most newest snapshot
func (so *SnapshotObject) ReadNewest(){
	snapshots := checkAvailableSnapshots(so.Dir)
	if len(snapshots) == 0{
		log.Info("Can't find available snapshots")
	} else {
		stat, err := os.Stat(snapshots[0])
		if err != nil {
			log.Fatal(err)
		}
		modtime := stat.ModTime()

		for _, fi := range snapshots {
			item, err := os.Stat(fi)
			if err != nil {
				log.Fatal(err)
			}

			if item.ModTime().After(modtime){
				modtime = item.ModTime()
			}
		}
	}
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
