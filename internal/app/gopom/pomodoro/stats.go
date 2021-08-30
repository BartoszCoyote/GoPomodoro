package pomodoro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type PomodoroStats struct {
	Enabled          bool
	DoneValue        int `json:"done"`
	InterruptedValue int `json:"failed"`
	RestsValue       int `json:"rests"`
}

func (ps *PomodoroStats) read() {
	if !ps.Enabled {
		return
	}
	f, err := os.OpenFile("/tmp/pomodoro-stats", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Cannot open stats file")
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f) // Error handling elided for brevity.
	if err != nil {
		log.Fatal("Cannot read file")
	}

	_ = json.Unmarshal(buf.Bytes(), &ps)
}
func (ps *PomodoroStats) save() {
	if !ps.Enabled {
		return
	}
	f, err := os.OpenFile("/tmp/pomodoro-stats", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Cannot open stats file")
	}
	defer f.Close()

	data, err := json.Marshal(ps)
	if err != nil {
		log.Fatal("Unable to encode json of stats.")
	}

	truncateWriteToFile(f, string(data))
}

func (ps *PomodoroStats) Done() {
	ps.read()
	ps.DoneValue += 1
	ps.save()
}

func (ps *PomodoroStats) Interrupted() {
	ps.read()
	ps.InterruptedValue += 1
	ps.save()
}

func (ps *PomodoroStats) Rest() {
	ps.read()
	ps.RestsValue += 1
	ps.save()
}

func truncateWriteToFile(f *os.File, text string) {
	err := f.Truncate(0)
	if err != nil {
		log.Fatal("Unable to truncate file for stats saving.")
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal("Unable to seek file for stats saving.")
	}
	_, err = fmt.Fprintf(f, text)
	if err != nil {
		log.Fatal(err)
	}
}
