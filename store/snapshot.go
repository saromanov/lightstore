// snapshot provides creating of shanpshots
package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"

	"github.com/op/go-logging"
	"github.com/saromanov/lightstore/proto"
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
	st    *Store
}

//NewSnapshotObject object provides initialization od new snapshot
func NewSnapshotObject(st *Store, path string) *SnapshotObject {
	so := new(SnapshotObject)
	so.Dir = "."
	so.st = st
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

// WriteNew is a temp method for writing of
// snapshots with entries
func (so *SnapshotObject) WriteNew(w io.Writer) error {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
	return nil
}

func (so *SnapshotObject) write(w io.Writer) error {
	kv := &protos.KVPair{}
	buf, err := kv.Marshal()
	if err != nil {
		return fmt.Errorf("unable to marshal data to snapshot: %v", err)
	}
	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("unable to write data to snapshot: %v", err)
	}
	return nil
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
