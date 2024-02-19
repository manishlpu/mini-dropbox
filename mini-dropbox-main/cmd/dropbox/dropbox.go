package main

import (
	"os"

	"github.com/manishlpu/assignment/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dropbox",
	Short: "Root command of the mini dropbox project",
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.ErrorLog("could not execute min-dropbox", err)
		os.Exit(1)
	}
}
