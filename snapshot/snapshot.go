// snapshot provides creating of shanpshots
package snapshot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("lightstore_log")

//Basic snapshot for all data in ligtstore

type snapshot interface {
	Write(object SnapshotObject)
}

type SnapshotObject struct {
	Crc32 string
	Data  string
	Dir   string
}

//NewSnapshotObject object provides initialization od new snapshot
func NewSnapshotObject(path string) *SnapshotObject {
	so := new(SnapshotObject)
	so.Dir = "."
	if path != "" {
		so.Dir = path
	}
	return so
}

//Write provides storing of new snapshot
func (so *SnapshotObject) Write() {
	b, err := json.Marshal(so)
	if err != nil {
		panic(err)
	}

	errwrite := ioutil.WriteFile(path.Join(so.Dir, "snap1.lssnapshot"), b, 0777)
	if errwrite != nil {
		panic(err)
	}

}

//Read provides reading snapshot and store data to lightstore
//if name is ""(empty), load more recently snapshot
func (so *SnapshotObject) Read(name string) {
	_, err := ioutil.ReadFile(path.Join(so.Dir, name))
	if err != nil {
		panic(fmt.Sprintf("Can't find snapshot with the name %s", name))
	}

	snapshots := checkAvailableSnapshots(so.Dir)
	if len(snapshots) == 0 {
		log.Info("Can't find available snapshots")
	} else {

	}
}

//Read newest provides reading most newest snapshot
func (so *SnapshotObject) ReadNewest() {
	snapshots := checkAvailableSnapshots(so.Dir)
	if len(snapshots) == 0 {
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

			if item.ModTime().After(modtime) {
				modtime = item.ModTime()
			}
		}
	}
}

// checkAvailableSnapshots Returns snapshot names
func checkAvailableSnapshots(dir string) []string {
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
