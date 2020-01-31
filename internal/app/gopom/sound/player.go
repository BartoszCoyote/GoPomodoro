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

func (p *Player) PlayLoop() {
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

func (p *Player) Play() {
	p.Stop()
	done := make(chan bool)
	fmt.Println("1 finish start")
	sound := p.sound.Streamer(0, p.sound.Len())
	speaker.Play(beep.Seq(sound, beep.Callback(func() {
		fmt.Println("3 finish done")
		done <- true
	})))
	fmt.Println("2 finish waiting")
	<-done
	fmt.Println("4 finish stop")
}

func (p *Player) Stop() {
	close(p.stopChan)
	close(p.loopChan)
}

func load_sound(soundFile string) *beep.Buffer {
	timer, err := os.Open(soundFile)
	if err != nil {
		fmt.Println("Fatal error reading sound file ")
	}

	streamer, format, err := vorbis.Decode(timer)
	if err != nil {
		fmt.Println("unable to decode ogg vorbis file")
	}
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Speaker initialization unsuccessful: ", err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return buffer
}
