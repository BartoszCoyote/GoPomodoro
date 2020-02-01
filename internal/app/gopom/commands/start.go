package commands

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/config"
	"github.com/cheggaaa/pb/v3"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
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

		timer, err := os.Open("./timer.mp3")
		if err != nil {
			fmt.Print("Fatal error reading sound file ")
		}

		streamer, format, err := mp3.Decode(timer)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		if err != nil {
			fmt.Print("unable to decode mp3 file")
		}
		defer streamer.Close()

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))

		// 25 minutes
		seconds := getTimers()

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
		<-done

		bar.Finish()
	},
}

func getTimers() int {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration config.Configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configuration.Pomodoro.Timer
}
