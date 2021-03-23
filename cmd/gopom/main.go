package main

import (
	"bufio"
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/commands"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
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

	go func(input chan rune) {
		for {
			reader := bufio.NewReader(os.Stdin)
			r, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println("error when reading rune from the stdin reader", err)
			}
			input <- r
		}
	}(pomodoro.StdinChan)

	commands.Execute()
}
