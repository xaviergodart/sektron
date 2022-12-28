package main

import (
	"fmt"
	"log"
	"os"
	"sektron/midi"
	"sektron/sequencer"
	"sektron/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	midi, err := midi.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()
	midi.Start()

	seq := sequencer.New(midi)

	p := tea.NewProgram(ui.New(seq))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
