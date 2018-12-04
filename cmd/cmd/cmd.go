package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:               "lightstore",
	Short:             "Command line interface for Lightstore",
	PersistentPreRunE: validateRootCmdArgs,
}

func validators(cmd *cobra.Command, args []string) error {
	if strings.HasPrefix(cmd.Use, "help ") { // No need to validate if it is help
		return nil
	}
	if sstDir == "" {
		return errors.New("--sst-dir not specified")
	}
	return nil
}

//Execute provides execution of the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
