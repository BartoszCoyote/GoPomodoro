package main

import (
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/commands"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/slack"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func main() {

	//TODO: introduce cleanup funcion/module that would contain all the logic related to `cleanup`
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		if viper.GetBool("ENABLE_SLACK_DND") {
			slack.EndDnd()
		}
		os.Exit(0)
	}()

	commands.Execute()
}
