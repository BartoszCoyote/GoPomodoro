package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start [taskName]",
	Short: "Start a task",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := getTaskName(args)
		pomodoro.NewPomodoro(taskName, 25*60, 5*60, 20*60, 4).Start()
	},
}

func getTaskName(args []string) string {
	var taskName = "task"
	if len(args) == 0 {
		fmt.Println("You haven't provided a task name. I will call it just a \"task\" for you.")
	} else {
		taskName = args[0]
	}
	return taskName
}
