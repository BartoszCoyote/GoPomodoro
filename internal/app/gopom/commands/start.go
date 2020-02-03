package commands

import (
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start task",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]
		pomodoro.NewPomodoro(taskName, 500, 500, 500, 4).Start()
	},
}
