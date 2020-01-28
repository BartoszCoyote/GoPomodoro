package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/sound"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"time"
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

		player := sound.NewPlayer("./timer_short.ogg")
		go player.Play()
		defer player.Stop()

		// 25 minutes
		seconds := 10
		//seconds := 25 * 60

		tmpl := `{{ "Working on task:" }} {{ string . "task" | green }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{string . "timer" | green}}`
		// start bar based on our template
		bar := pb.ProgressBarTemplate(tmpl).Start(seconds)
		// set values for string elements
		bar.Set("task", taskName).
			Set("timer", fmtTimer(0))

		for i := 0; i < seconds; i++ {
			bar.Increment()
			time.Sleep(1 * time.Second)
			bar.Set("timer", fmtTimer(i))

		}

		bar.Finish()
	},
}
