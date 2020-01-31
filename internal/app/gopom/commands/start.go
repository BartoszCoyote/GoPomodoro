package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/task"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

func fmtTimer(t int) string {
	m := t / 60
	s := t - (m * 60)
	return fmt.Sprintf("%02d:%02d", m, s)
}

// Beep documentation - https://github.com/faiface/beep/wiki/Hello,-Beep!
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start task",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]

		fmt.Print("started task :", taskName)

		subtask := task.NewSubtask(taskName, 3, "./timer_short.ogg", "./timer_short2.ogg")
		subtask.Work()
	},
}
