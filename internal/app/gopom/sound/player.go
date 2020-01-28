package sound

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"os"
	"time"
)

type Player struct {
	sound    *beep.Buffer
	loopChan chan bool
	stopChan chan struct{}
}

func NewPlayer(soundFile string) *Player {
	return new(Player).Init(soundFile)
}

func (p *Player) Init(soundFile string) *Player {
	p.loopChan = make(chan bool)
	p.stopChan = make(chan struct{})
	p.sound = load_sound(soundFile)
	return p
}

func (p *Player) Play() {
	for {
		sound := p.sound.Streamer(0, p.sound.Len())
		speaker.Play(beep.Seq(sound, beep.Callback(func() {
			p.loopChan <- true
		})))
		select {
		case <-p.loopChan:
			continue
		case <-p.stopChan:
			break
		}
	}
}

func (p *Player) Stop() {
	close(p.stopChan)
	close(p.loopChan)
}

func load_sound(soundFile string) *beep.Buffer {
	timer, err := os.Open(soundFile)
	if err != nil {
		fmt.Print("Fatal error reading sound file ")
	}

	streamer, format, err := vorbis.Decode(timer)
	if err != nil {
		fmt.Print("unable to decode mp3 file")
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return buffer
}
