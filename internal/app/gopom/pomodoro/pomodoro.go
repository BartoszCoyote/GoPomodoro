package pomodoro

import (
	"fmt"
	"time"

	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/sound"
	"github.com/cheggaaa/pb/v3"
	"github.com/looplab/fsm"
)

const (
	INITIALIZED_STATE           string = "init"
	WORK_STATE                  string = "working"
	WORK_COUNT_EVALUATION_STATE string = "counting_work"
	REST_STATE                  string = "resting"
	LONG_REST_STATE             string = "long_resting"

	WORK_STARTED_EVENT        string = "started"
	WORK_FINISHED_EVENT       string = "work_finished"
	REST_FINISHED_EVENT       string = "rest_finished"
	MORE_WORK_NEEDED_EVENT    string = "more_work_needed"
	NO_MORE_WORK_NEEDED_EVENT string = "no_more_work_needed"
	WORK_RESTARTED_EVENT      string = "restarted"
)

type Pomodoro struct {
	taskName         string
	workDuration     int
	restDuration     int
	longRestDuration int
	cycles           int
	maxCycles        int
	stateMachine     *fsm.FSM
}

type Subtask struct {
	duration    int
	name        string
	workSound   *sound.Player
	finishSound *sound.Player
	progress    *pb.ProgressBar
}

func NewPomodoro(taskName string, workDuration int, restDuration int, longRestDuration int, maxCycles int) *Pomodoro {
	return &Pomodoro{
		taskName,
		workDuration,
		restDuration,
		longRestDuration,
		0,
		maxCycles,
		initStateMachine(),
	}
}

func NewSubtask(name string, duration int, workSound string, finishSound string) *Subtask {
	barTemplate := `{{ string . "task" | green }} {{ bar . "▇" "▇" (cycle . "▂" "▃" "▅" "▆" "▅" "▃" "▂" ) "_" "▇"}} {{string . "timer" | green}}`
	return &Subtask{
		duration,
		name,
		sound.NewPlayer(workSound),
		sound.NewPlayer(finishSound),
		pb.ProgressBarTemplate(barTemplate).
			Start(duration).
			Set("task", name).
			Set("timer", fmtTimer(0)),
	}
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

func (p *Pomodoro) Start() {
	stateHandler := make(map[string]func() string)

	stateHandler[INITIALIZED_STATE] = p.init
	stateHandler[WORK_STATE] = p.work
	stateHandler[WORK_COUNT_EVALUATION_STATE] = p.evaluateWorkCount
	stateHandler[REST_STATE] = p.rest
	stateHandler[LONG_REST_STATE] = p.longRest

	for {
		handler := stateHandler[p.stateMachine.Current()]
		event := handler()
		err := p.stateMachine.Event(event)

		if err != nil {
			fmt.Println("State transition unsuccessful: ", err)
		}
	}
}

func initStateMachine() *fsm.FSM {
	return fsm.NewFSM(
		INITIALIZED_STATE,
		fsm.Events{
			{Src: []string{INITIALIZED_STATE}, Name: WORK_STARTED_EVENT, Dst: WORK_STATE},
			{Src: []string{WORK_STATE}, Name: WORK_FINISHED_EVENT, Dst: WORK_COUNT_EVALUATION_STATE},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: MORE_WORK_NEEDED_EVENT, Dst: REST_STATE},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: NO_MORE_WORK_NEEDED_EVENT, Dst: LONG_REST_STATE},
			{Src: []string{REST_STATE}, Name: REST_FINISHED_EVENT, Dst: WORK_STATE},
			{Src: []string{LONG_REST_STATE}, Name: WORK_RESTARTED_EVENT, Dst: INITIALIZED_STATE},
		},
		fsm.Callbacks{
			//"enter_state": func(event *fsm.Event) {
			//	fmt.Printf("from %s to %s\n", event.Src, event.Dst)
			//},
		},
	)
}

func (p *Pomodoro) init() string {
	taskStartupName := "Starting work on " + p.taskName
	subtask := NewSubtask(taskStartupName, 2, "/beep.mp3", "/placeholder.mp3")
	subtask.Work()

	return WORK_STARTED_EVENT
}

func (p *Pomodoro) work() string {
	workName := "Working on " + p.taskName
	subtask := NewSubtask(workName, p.workDuration, "/timer.mp3", "/finish.mp3")
	subtask.Work()

	p.cycles++

	return WORK_FINISHED_EVENT
}

func (p *Pomodoro) evaluateWorkCount() string {
	if p.cycles >= p.maxCycles {
		fmt.Println("You have finished working on task ", p.taskName)
		return NO_MORE_WORK_NEEDED_EVENT
	}
	return MORE_WORK_NEEDED_EVENT
}

func (p *Pomodoro) rest() string {
	subtask := NewSubtask("Resting...", p.restDuration, "/placeholder.mp3", "/finish.mp3")
	subtask.Work()

	return REST_FINISHED_EVENT
}

func (p *Pomodoro) longRest() string {
	subtask := NewSubtask("Long rest...", p.longRestDuration, "/placeholder.mp3", "/finish.mp3")
	subtask.Work()
	p.cycles = 0

	return WORK_RESTARTED_EVENT
}

func fmtTimer(t int) string {
	m := t / 60
	s := t - (m * 60)
	return fmt.Sprintf("%02d:%02d", m, s)
}
