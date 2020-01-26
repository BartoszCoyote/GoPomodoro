package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the go pomodoro app version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Local development version of go pomodoro.")
	},
}
