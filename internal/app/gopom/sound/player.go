package sound

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"sync"
	"time"
)

var (
	mu                 sync.Mutex
	speakerInitialized bool
)

type Player struct {
	sound    beep.StreamSeekCloser
	stopChan chan struct{}
}

func NewPlayer(soundFile string) *Player {
	return &Player{
		loadSound(soundFile),
		make(chan struct{}),
	}
}

func (p *Player) PlayLoop() {
	controlledSound := &beep.Ctrl{Streamer: beep.Loop(-1, p.sound), Paused: false}
	speaker.Play(controlledSound)

	<-p.stopChan
	speaker.Lock()
	controlledSound.Paused = true
	speaker.Unlock()
}

func (p *Player) Stop() {
	close(p.stopChan)
}

func (p *Player) Play() {
	done := make(chan struct{})
	speaker.Play(beep.Seq(p.sound, beep.Callback(func() {
		close(done)
	})))
	<-done
}

func loadSound(soundFile string) beep.StreamSeekCloser {
	timer, err := os.Open(soundFile)
	if err != nil {
		fmt.Println("Fatal error reading sound file ")
	}

	streamer, format, err := mp3.Decode(timer)
	if err != nil {
		fmt.Println("unable to decode mp3 file")
	}
	initSpeaker(format)

	return streamer
}

func initSpeaker(format beep.Format) {
	mu.Lock()
	defer mu.Unlock()

	if !speakerInitialized {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			fmt.Println("Speaker initialization unsuccessful: ", err)
		}
		speakerInitialized = true
	}
}
