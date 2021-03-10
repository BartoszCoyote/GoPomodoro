package main

import (
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/commands"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/slack"
	"os"
	"os/signal"
)

func main() {

	//TODO: introduce cleanup funcion/module that would contain all the logic related to `cleanup`
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		slack.EndDnd()
		os.Exit(0)
	}()

	commands.Execute()
}
