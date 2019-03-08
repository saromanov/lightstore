package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/saromanov/lightstore/logging"
	"github.com/saromanov/lightstore/proto"
)

// Snapshot provides creating of the new snapshoе
type Snapshot struct {
	Crc32 string
	Data  string
	Dir   string
	st    *Store
}

//NewSnapshot object provides initialization od new snapshot
func NewSnapshot(st *Store, path string) *Snapshot {
	dir := "."
	if path != "" {
		dir = path
	}
	return &Snapshot{
		st:  st,
		Dir: dir,
	}
}

// Write is a method for writing of
// snapshots with entries
func (so *Snapshot) Write(w io.Writer) error {
	buf := new(bytes.Buffer)
	txn := so.st.NewTransaction(false)
	it, err := txn.NewIterator(IteratorOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to make iterator")
	}
	for it.First(); it.Valid(); it.Next() {
		itm := it.Item()
		data := &protos.KVPair{
			Key:   itm.Key(),
			Value: itm.ValueData(),
		}
		err := binary.Write(buf, binary.LittleEndian, data)
		if err != nil {
			return errors.Wrap(err, "unable to write data")
		}
		fmt.Printf("% x", buf.Bytes())
	}
	return nil
}

//Read provides reading snapshot and store data to lightstore
//if name is ""(empty), load more recently snapshot
func (so *Snapshot) Read(name string) {
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
func (so *Snapshot) ReadNewest() {
	snapshots := checkAvailableSnapshots(so.Dir)
	if len(snapshots) == 0 {
		log.Info("Can't find available snapshots")
	} else {
		stat, err := os.Stat(snapshots[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		modtime := stat.ModTime()

		for _, fi := range snapshots {
			item, err := os.Stat(fi)
			if err != nil {
				log.Fatal(err.Error())
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
