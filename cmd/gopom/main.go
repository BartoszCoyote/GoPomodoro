package main

import (
	"bufio"
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/commands"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/pomodoro"
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
