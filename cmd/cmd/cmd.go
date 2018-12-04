package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "lightstore",
	Short:             "Command line interface for Lightstore",
	PersistentPreRunE: validator,
}

func validator(cmd *cobra.Command, args []string) error {
	if strings.HasPrefix(cmd.Use, "help ") {
		return nil
	}
	return nil
}

//Execute provides execution of the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
