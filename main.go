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
	configFile := flag.String("config", "config.json", "config file to load or create")
	keyboard := flag.String("keyboard", "", "keyboard layout (qwerty, qwerty-mac, azerty, azerty-mac)")
	patternsFile := flag.String("patterns", "patterns.json", "patterns file to load or create")
	flag.Parse()

	midi, err := midi.New()
	if err != nil {
		log.Fatal(err)
	}
	defer midi.Close()

	config := filesystem.NewConfiguration(*configFile, *keyboard)
	bank := filesystem.NewBank(*patternsFile)

	seq := sequencer.New(midi, bank)

	p := tea.NewProgram(ui.New(config, seq))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
