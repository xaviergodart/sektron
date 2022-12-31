package main

import (
	"log"
	"sektron/midi"
	"sektron/sequencer"
	"sektron/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	midi, err := midi.NewMidi()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	seq := sequencer.New(midi)

	p := tea.NewProgram(ui.New(seq))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
