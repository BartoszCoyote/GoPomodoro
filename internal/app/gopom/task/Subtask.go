package task

import (
	"fmt"
	sound "github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/sound"
	"github.com/cheggaaa/pb/v3"
	"time"
)

type Subtask struct {
	duration    int
	name        string
	workSound   *sound.Player
	finishSound *sound.Player
	progress    *pb.ProgressBar
}

func NewSubtask(name string, duration int, workSound string, finishSound string) *Subtask {
	return new(Subtask).Init(name, duration, workSound, finishSound)
}

func (s *Subtask) Init(name string, duration int, workSound string, finishSound string) *Subtask {
	s.name = name
	s.duration = duration

	s.workSound = sound.NewPlayer(workSound)
	s.finishSound = sound.NewPlayer(finishSound)

	tmpl := `{{ string . "task" | green }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{string . "timer" | green}}`
	s.progress = pb.ProgressBarTemplate(tmpl).Start(s.duration)
	s.progress.
		Set("task", s.name).
		Set("timer", fmtTimer(0))
	return s
}

func (s *Subtask) Work() {
	go s.workSound.PlayLoop()

	for i := 0; i < s.duration; i++ {
		s.progress.Increment()
		time.Sleep(1 * time.Second)
		s.progress.Set("timer", fmtTimer(i))
	}

	s.workSound.Stop()
	s.progress.Finish()
	s.finishSound.Play()
}

func fmtTimer(t int) string {
	m := t / 60
	s := t - (m * 60)
	return fmt.Sprintf("%02d:%02d", m, s)
}
