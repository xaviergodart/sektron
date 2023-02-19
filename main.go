package main

import (
	"log"
	"sektron/filesystem"
	"sektron/midi"
	"sektron/sequencer"
	"sektron/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	// load default saved pattern
	bank := filesystem.NewBank("patterns.json")

	seq := sequencer.New(midi, bank)

	p := tea.NewProgram(ui.New(seq))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
