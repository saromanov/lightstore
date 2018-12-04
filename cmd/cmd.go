package cmd

import (
	"flag"
	"log"
	"strings"

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

func parseCommands() {
	flag.Parse()
	if *put != "" {
		pair := strings.Split(*put, "")
		if len(pair) != 2 {
			log.Fatal("unable to get key value pair from 'put' command. Should be 'key value'")
		}
	}
}
func main() {
	initLightStore()
	parseCommands()
}
