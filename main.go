package main

import (
	"flag"
	"log"
	"sektron/filesystem"
	"sektron/midi"
	"sektron/sequencer"
	"sektron/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	patternsFile := flag.String("patterns", "patterns.json", "patterns file to load or create")
	flag.Parse()

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	bank := filesystem.NewBank(*patternsFile)

	seq := sequencer.New(midi, bank)

	p := tea.NewProgram(ui.New(seq))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
