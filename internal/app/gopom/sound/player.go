package sound

import (
	"fmt"
	"github.com/faiface/beep/effects"
	"log"
	"sync"
	"time"

	// required for statik file system to enable embedded resources
	_ "github.com/BartoszCoyote/GoPomodoro/internal/app/gopom/sound/soundpack"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/rakyll/statik/fs"
)

var (
	mu                 sync.Mutex
	speakerInitialized bool
)

type Player struct {
	sound    *beep.StreamSeekCloser
	volume   float64
	mute     bool
	stopChan chan struct{}
}

// Volume should be in range of -8..2 where 0 is unchanged system volume, -8 is barely audible, 2 is loud but not yet
// distorted much. Setting volume higher than 2 may be unpleasant
func NewPlayer(soundFile string, soundVolume float64) *Player {
	sound, _ := loadSound(soundFile)

	return &Player{
		sound,
		soundVolume,
		soundVolume <= -8,
		make(chan struct{}),
	}
}

func (p *Player) PlayLoop() {
	volume := p.withVolumeEffect(beep.Loop(-1, *p.sound))
	speaker.Play(volume)

	<-p.stopChan
	speaker.Lock()
	volume.Streamer = beep.Loop(0, *p.sound)
	speaker.Unlock()
}

func (p *Player) Stop() {
	close(p.stopChan)
}

func (p *Player) Play() {
	done := make(chan struct{})
	speaker.Play(beep.Seq(p.withVolumeEffect(beep.Loop(1, *p.sound)), beep.Callback(func() {
		close(done)
	})))
	<-done
}

func loadSound(soundFile string) (*beep.StreamSeekCloser, beep.Format) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	timer, err := statikFS.Open(soundFile)
	if err != nil {
		fmt.Println("Fatal error reading sound file ")
	}

	streamer, format, err := mp3.Decode(timer)
	if err != nil {
		fmt.Println("unable to decode mp3 file")
	}
	initSpeaker(format)

	return &streamer, format
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

func (p *Player) withVolumeEffect(streamer beep.Streamer) *effects.Volume {
	return &effects.Volume{Streamer: streamer, Base: 2, Silent: p.mute, Volume: p.volume}
}
