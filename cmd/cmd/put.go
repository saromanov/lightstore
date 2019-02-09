package cmd

import (
	"log"

	"github.com/saromanov/lightstore/store"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Putting of the key value pair",
	Long:  `Putting of the key value pair on format 'key value'`,
	Run:   put,
}

func put(cmd *cobra.Command, args []string) {
	ls, err := store.Open(nil)
	if err != nil {
		panic(err)
	}
	if len(args) != 2 {
		log.Fatal("put command should contains 2 arguments")
	}
	key := args[0]
	value := args[1]
	err := ls.Write(func(txn *store.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return txn.Commit()
	})
	if err != nil {
		log.Fatalf("unable to write data: %v", err)
	}
}