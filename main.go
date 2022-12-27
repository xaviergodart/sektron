package main

import (
    "fmt"
    "log"
    "os"
    "sektron/ui"

    tea "github.com/charmbracelet/bubbletea"
    "gitlab.com/gomidi/midi/v2"
    _ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

func main() {
    defer midi.CloseDriver()

    drivers := midi.GetOutPorts()
    if len(drivers) == 0 {
        log.Fatal("No midi drivers")
    }
    midiSend, _ := midi.SendTo(drivers[0])

    p := tea.NewProgram(ui.New(midiSend))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
