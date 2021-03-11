package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start [taskName]",
	Short: "Start a task",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := getTaskName(args)
		workDuration := viper.GetInt("WORK_DURATION_MINUTES")
		restDuration := viper.GetInt("REST_DURATION_MINUTES")
		longRestDuration := viper.GetInt("LONG_REST_DURATION_MINUTES")
		maxCycles := viper.GetInt("MAX_CYCLES")
		pomodoro.NewPomodoro(taskName, workDuration*60, restDuration*60, longRestDuration*60, maxCycles).Start()
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
