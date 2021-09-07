package pomodoro

import (
	"fmt"
	"github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/slack"
	"github.com/spf13/viper"
	"log"
	"os"
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
	WORK_CONTINUE_PROMPT        string = "waiting_for_user"

	WORK_STARTED_EVENT        string = "started"
	WORK_FINISHED_EVENT       string = "work_finished"
	REST_FINISHED_EVENT       string = "rest_finished"
	MORE_WORK_NEEDED_EVENT    string = "more_work_needed"
	NO_MORE_WORK_NEEDED_EVENT string = "no_more_work_needed"
	WORK_RESTARTED_EVENT      string = "restarted"
	WORK_RESUMED_EVENT        string = "resumed"
	WORK_INTERRUPTED_EVENT    string = "interrupted"
)

type PomodoroSettings struct {
	TaskName                string
	WorkDurationMinutes     int
	RestDurationMinutes     int
	LongRestDurationMinutes int
	Cycles                  int
	WorkSoundVolume         float64
	FinishSoundVolume       float64
	TimerToFileEnabled      bool
	StatsToFileEnabled      bool
}

type Pomodoro struct {
	taskName                string
	workDurationMinutes     int
	restDurationMinutes     int
	longRestDurationMinutes int
	cycles                  int
	maxCycles               int
	workSoundVolume         float64
	finishSoundVolume       float64
	stateMachine            *fsm.FSM
	enableFileTimer         bool
	stats                   *PomodoroStats
}

type Subtask struct {
	name            string
	workSound       *sound.Player
	finishSound     *sound.Player
	progress        *pb.ProgressBar
	echoTimerToFile bool
}

var StdinChan = make(chan rune)

func NewPomodoro(pomodoroSettings *PomodoroSettings) *Pomodoro {

	stats := &PomodoroStats{}
	stats.Enabled = pomodoroSettings.StatsToFileEnabled

	return &Pomodoro{
		pomodoroSettings.TaskName,
		pomodoroSettings.WorkDurationMinutes,
		pomodoroSettings.RestDurationMinutes,
		pomodoroSettings.LongRestDurationMinutes,
		0,
		pomodoroSettings.Cycles,
		pomodoroSettings.WorkSoundVolume,
		pomodoroSettings.FinishSoundVolume,
		initStateMachine(),
		pomodoroSettings.TimerToFileEnabled,
		stats,
	}
}

func newSubtask(name string, durationInSeconds int, workSound string, workSoundVolume float64, finishSound string, finishSoundVolume float64, fileTimer bool) *Subtask {
	barTemplate := `{{ string . "task" | green }} {{ bar . "▇" "▇" (cycle . "▂" "▃" "▅" "▆" "▅" "▃" "▂" ) "_" "▇"}} {{string . "timer" | green}}`
	totalSubtaskTime := int64(durationInSeconds) * time.Second.Milliseconds()
	return &Subtask{
		name,
		sound.NewPlayer(workSound, workSoundVolume),
		sound.NewPlayer(finishSound, finishSoundVolume),
		pb.ProgressBarTemplate(barTemplate).
			Start64(totalSubtaskTime).
			Set("task", name).
			Set("timer", fmtTimer(0, totalSubtaskTime)),
		fileTimer,
	}
}

func (s *Subtask) work() bool {
	go s.workSound.PlayLoop()

	// stop the work sound and progress regardless of the outcome of the task - finished, partially finished
	defer s.workSound.Stop()
	defer s.progress.Finish()
	defer s.finishSound.Play()

	var f *os.File
	var err error
	if s.echoTimerToFile {
		f, err = os.OpenFile("/tmp/pomodoro-timer", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	if f != nil {
		defer f.Close()
	}

	for {
		select {
		case <-StdinChan:
			return false
		case <-time.After(1000 * time.Millisecond):
			currentDuration := time.Now().Sub(s.progress.StartTime()).Milliseconds()
			s.progress.SetCurrent(currentDuration)

			if s.progress.Total() < currentDuration {
				return true
			}
			s.progress.Set("timer", fmtTimer(currentDuration, s.progress.Total()))

			if f != nil {
				truncateWriteToFile(f, fmtTimer(currentDuration, s.progress.Total()))
			}
		}
	}
}

func (p *Pomodoro) Start() {
	stateHandler := make(map[string]func() string)

	stateHandler[INITIALIZED_STATE] = p.init
	stateHandler[WORK_STATE] = p.work
	stateHandler[WORK_COUNT_EVALUATION_STATE] = p.evaluateWorkCount
	stateHandler[REST_STATE] = p.rest
	stateHandler[LONG_REST_STATE] = p.longRest
	stateHandler[WORK_CONTINUE_PROMPT] = p.waitForUser

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
			{Src: []string{WORK_STATE}, Name: WORK_INTERRUPTED_EVENT, Dst: WORK_CONTINUE_PROMPT},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: MORE_WORK_NEEDED_EVENT, Dst: REST_STATE},
			{Src: []string{WORK_COUNT_EVALUATION_STATE}, Name: NO_MORE_WORK_NEEDED_EVENT, Dst: LONG_REST_STATE},
			{Src: []string{REST_STATE}, Name: REST_FINISHED_EVENT, Dst: WORK_CONTINUE_PROMPT},
			{Src: []string{WORK_CONTINUE_PROMPT}, Name: WORK_RESUMED_EVENT, Dst: WORK_STATE},
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
	subtask := newSubtask(taskStartupName, 2, "/beep.mp3", p.workSoundVolume, "/placeholder.mp3", p.finishSoundVolume, p.enableFileTimer)
	subtask.work()

	return WORK_STARTED_EVENT
}

func (p *Pomodoro) work() string {
	slackDndEnabled := viper.GetBool("ENABLE_SLACK_DND")
	if slackDndEnabled {
		slack.SetDnd(p.workDurationMinutes)
	}

	workName := "Working on " + p.taskName
	subtask := newSubtask(workName, p.workDurationMinutes*60, "/timer.mp3", p.workSoundVolume, "/finish.mp3", p.finishSoundVolume, p.enableFileTimer)
	workFinished := subtask.work()

	if slackDndEnabled {
		slack.EndDnd()
	}

	if workFinished {
		p.cycles++
		p.stats.Done()
		return WORK_FINISHED_EVENT
	} else {
		fmt.Println("Pomodoro was interrupted.")
		p.stats.Interrupted()
		return WORK_INTERRUPTED_EVENT
	}
}

func (p *Pomodoro) evaluateWorkCount() string {
	if p.cycles >= p.maxCycles {
		fmt.Println("You have finished working on task ", p.taskName)
		return NO_MORE_WORK_NEEDED_EVENT
	}
	return MORE_WORK_NEEDED_EVENT
}

func (p *Pomodoro) rest() string {
	subtask := newSubtask("Resting...", p.restDurationMinutes*60, "/placeholder.mp3", p.workSoundVolume, "/finish.mp3", p.finishSoundVolume, p.enableFileTimer)
	subtask.work()

	p.stats.Rest()
	return REST_FINISHED_EVENT
}

func (p *Pomodoro) longRest() string {
	subtask := newSubtask("Long rest...", p.longRestDurationMinutes*60, "/placeholder.mp3", p.workSoundVolume, "/finish.mp3", p.finishSoundVolume, p.enableFileTimer)
	subtask.work()
	p.cycles = 0

	p.stats.Rest()
	return WORK_RESTARTED_EVENT
}

func (p *Pomodoro) waitForUser() string {
	waitForUser := viper.GetBool("ENABLE_WORK_CONTINUE")
	if waitForUser {
		fmt.Println("Press Enter to continue...")
		<-StdinChan
	}

	return WORK_RESUMED_EVENT
}

func fmtTimer(currentMs int64, totalMs int64) string {
	currentMin := currentMs / 1000 / 60
	currentSec := currentMs/1000 - (currentMin * 60)
	totalMin := totalMs / 1000 / 60
	totalSec := totalMs/1000 - (totalMin * 60)
	return fmt.Sprintf("%02d:%02d / %02d:%02d", currentMin, currentSec, totalMin, totalSec)
}
