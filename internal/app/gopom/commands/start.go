package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
	"github.com/spf13/cobra"
)

func init() {
	startCmd.Flags().Float64VarP(&workVolume, "work-volume", "w", 80, "Sets volume of sound when working. Values 0..100.")
	startCmd.Flags().Float64VarP(&finishVolume, "finish-volume", "f", 80, "Sets volume of sound when finished working. Values 0..100.")
	startCmd.Flags().BoolVar(&muteWorkSounds, "mute-work", false, "Disables sound played when working.")
	startCmd.Flags().BoolVar(&muteFinishSounds, "mute-finish", false, "Disables sound played when finished.")
	rootCmd.AddCommand(startCmd)
}

var workVolume float64
var finishVolume float64

var muteWorkSounds bool
var muteFinishSounds bool

var startCmd = &cobra.Command{
	Use:   "start [taskName]",
	Short: "Start a task",
	Run: func(cmd *cobra.Command, args []string) {
		if muteWorkSounds {
			workVolume = 0
		}
		if muteFinishSounds {
			finishVolume = 0
		}

		//converting from human readable 0..100 range to Players -8..2 range
		workVolume = workVolume/10 - 8
		finishVolume = finishVolume/10 - 8

		taskName := getTaskName(args)
		pomodoro.NewPomodoro(&pomodoro.PomodoroSettings{
			TaskName:          taskName,
			WorkDuration:      25 * 60,
			RestDuration:      5 * 60,
			LongRestDuration:  20 * 60,
			Cycles:            4,
			WorkSoundVolume:   workVolume,
			FinishSoundVolume: finishVolume,
		}).Start()
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
