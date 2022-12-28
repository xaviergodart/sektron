package main

import (
	"fmt"
	"log"
	"os"
	"sektron/midi"
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

	p := tea.NewProgram(ui.New(midi))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
