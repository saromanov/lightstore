package cmd

import (
	"github.com/spf13/cobra"
	"github.com/saromanov/lightstore/store"

)
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Putting of the key value pair",
	Long: `Putting of the key value pair on format 'key value'`,
	Run: func(cmd *cobra.Command, args []string) {
		ls = store.Open(nil)
	},
}

func put(cmd *cobra.Command, args []string) {
	ls := store.Open(nil)
}