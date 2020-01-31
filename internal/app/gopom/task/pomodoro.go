package task

import (
	"fmt"
	"github.com/looplab/fsm"
)

const (
	NOT_INITIALIZED_STATE       string = "not_init"
	INITIALIZED_STATE           string = "init"
	WORK_STATE                  string = "working"
	WORK_COUNT_EVALUATION_STATE string = "counting_work"
	REST_STATE                  string = "resting"
	LONG_REST_STATE             string = "long_resting"

	INITIALIZED_EVENT         string = "initialized"
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

func NewPomodoro(taskName string, workDuration int, restDuration int, longRestDuration int, maxCycles int) *Pomodoro {
	return new(Pomodoro).Init(taskName, workDuration, restDuration, longRestDuration, maxCycles)
}

func (p *Pomodoro) Init(taskName string, workDuration int, restDuration int, longRestDuration int, maxCycles int) *Pomodoro {
	p.taskName = taskName
	p.workDuration = workDuration
	p.restDuration = restDuration
	p.longRestDuration = longRestDuration
	p.maxCycles = maxCycles
	p.cycles = 0
	p.init_state_machine()

	err := p.stateMachine.Event(INITIALIZED_EVENT)
	consumeFsmEventTransitionError(err)

	return p
}

func (p *Pomodoro) init_state_machine() {
	p.stateMachine = fsm.NewFSM(
		NOT_INITIALIZED_STATE,
		fsm.Events{
			{Src: []string{NOT_INITIALIZED_STATE}, Name: INITIALIZED_EVENT, Dst: INITIALIZED_STATE},
			{Src: []string{INITIALIZED_STATE}, Name: WORK_STARTED_EVENT, Dst: WORK_STATE},
			{Src: []string{WORK_STATE}, Name: WORK_FINISHED_EVENT, Dst: WORK_COUNT_EVALUATION_STATE},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: MORE_WORK_NEEDED_EVENT, Dst: REST_STATE},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: NO_MORE_WORK_NEEDED_EVENT, Dst: LONG_REST_STATE},
			{Src: []string{REST_STATE}, Name: REST_FINISHED_EVENT, Dst: WORK_STATE},
			{Src: []string{LONG_REST_STATE}, Name: WORK_RESTARTED_EVENT, Dst: INITIALIZED_STATE},
		},
		fsm.Callbacks{
			INITIALIZED_STATE:               p.init,
			WORK_STATE:                      p.work,
			WORK_COUNT_EVALUATION_STATE:     p.evaluateWorkCount,
			REST_STATE:                      p.rest,
			LONG_REST_STATE:                 p.longRest,
			"after_" + WORK_RESTARTED_EVENT: p.restart,
		},
	)
}

func (p *Pomodoro) init(e *fsm.Event) {
	taskStartupName := "Starting work on " + p.taskName
	subtask := NewSubtask(taskStartupName, 2, "./beep.mp3", "./placeholder.mp3")
	subtask.Work()

	fmt.Println("1")
	fmt.Println(e.FSM.Current())
	err := e.FSM.Event(WORK_STARTED_EVENT)
	fmt.Println("2")
	fmt.Println(e.FSM.Current())
	consumeFsmEventTransitionError(err)
}

func (p *Pomodoro) work(e *fsm.Event) {
	workName := "Working on " + p.taskName
	subtask := NewSubtask(workName, p.workDuration, "./timer.mp3", "./finish.mp3")
	subtask.Work()

	p.cycles++

	err := e.FSM.Event(WORK_FINISHED_EVENT)
	consumeFsmEventTransitionError(err)
}

func (p *Pomodoro) evaluateWorkCount(e *fsm.Event) {
	if p.cycles >= p.maxCycles {
		fmt.Println("You have finished working on task ", p.taskName)
		err := e.FSM.Event(NO_MORE_WORK_NEEDED_EVENT)
		consumeFsmEventTransitionError(err)
	} else {
		err := e.FSM.Event(MORE_WORK_NEEDED_EVENT)
		consumeFsmEventTransitionError(err)
	}
}

func (p *Pomodoro) rest(e *fsm.Event) {
	subtask := NewSubtask("Resting...", p.restDuration, "./placeholder.mp3", "./finish.mp3")
	subtask.Work()

	err := e.FSM.Event(REST_FINISHED_EVENT)
	consumeFsmEventTransitionError(err)
}

func (p *Pomodoro) longRest(e *fsm.Event) {
	subtask := NewSubtask("Long rest...", p.longRestDuration, "./placeholder.mp3", "./finish.mp3")
	subtask.Work()

	err := e.FSM.Event(WORK_RESTARTED_EVENT)
	consumeFsmEventTransitionError(err)
}

func (p *Pomodoro) restart(e *fsm.Event) {
	p.cycles = 0

	fmt.Println("Started new Pomodoro!")
}

func consumeFsmEventTransitionError(err error) {
	if err != nil {
		fmt.Println("State machine transition error: ", err)
	}
}
