package cmd

import (
	"flag"

	"github.com/saromanov/lightstore/store"
)

var (
	put     = flag.String("put", "", "put of the key value pair")
	get     = flag.String("get", "", "get of the value by the key")
	backup  = flag.String("backup", "", "backup of the data")
	restore = flag.String("restore", "", "restore of the data")
)

var ls *store.Lightstore

func initLightStore() {
	ls = store.Open(nil)
}
func main() {
	initLightStore()
}
