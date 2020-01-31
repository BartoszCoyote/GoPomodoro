package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/task"
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

		fmt.Print("started task :", taskName)

		subtask := task.NewSubtask(taskName, 3, "./timer.mp3", "./finish.mp3")
		subtask.Work()
	},
}
