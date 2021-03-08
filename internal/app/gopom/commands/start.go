package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
	"github.com/spf13/cobra"
)

func init() {
	startCmd.Flags().IntVarP(&workVolume, "work-volume", "w", 80, "Sets volume of sound when working. Values 0..100.")
	startCmd.Flags().IntVarP(&finishVolume, "finish-volume", "f", 80, "Sets volume of sound when finished working. Values 0..100.")
	startCmd.Flags().BoolVar(&muteWorkSounds, "mute-work", false, "Disables sound played when working.")
	startCmd.Flags().BoolVar(&muteFinishSounds, "mute-finish", false, "Disables sound played when finished.")
	rootCmd.AddCommand(startCmd)
}

var workVolume int
var finishVolume int

var muteWorkSounds bool
var muteFinishSounds bool

var startCmd = &cobra.Command{
	Use:   "start [taskName]",
	Short: "Start a task",
	Run: func(cmd *cobra.Command, args []string) {
		if workVolume < 0 {
			workVolume = 0
		}
		if finishVolume < 0 {
			finishVolume = 0
		}
		if workVolume > 100 {
			workVolume = 100
		}
		if finishVolume > 100 {
			finishVolume = 100
		}

		if muteWorkSounds {
			workVolume = 0
		}
		if muteFinishSounds {
			finishVolume = 0
		}

		//converting from human readable 0..100 range to Players -8..2 range
		internalWorkVolume := float64(workVolume)/10 - 8
		internalFinishVolume := float64(finishVolume)/10 - 8

		taskName := getTaskName(args)
		pomodoro.NewPomodoro(&pomodoro.PomodoroSettings{
			TaskName:          taskName,
			WorkDuration:      25 * 60,
			RestDuration:      5 * 60,
			LongRestDuration:  20 * 60,
			Cycles:            4,
			WorkSoundVolume:   internalWorkVolume,
			FinishSoundVolume: internalFinishVolume,
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
