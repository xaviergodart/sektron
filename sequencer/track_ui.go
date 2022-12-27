package sequencer

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"gitlab.com/gomidi/midi/v2"
)

var (
	stepStyle = lipgloss.NewStyle().
			Width(6).
			Height(3).
			Margin(1).
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))
	stepCurrentStyle = lipgloss.NewStyle().
				Width(6).
				Height(3).
				Margin(1).
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFFFFF"))
)

func (t track) View(pulse int) string {
	var steps []string
	for i, step := range t.steps {
		if i == pulse/(pulsesPerQuarterNote/stepsPerQuarterNote) {
			if !step.triggered {
				msg := midi.NoteOn(0, step.note, 120)
				t.sendMidi(msg)
				step.triggered = true
			}
			steps = append(steps, stepCurrentStyle.Render(strconv.Itoa(i+1)))
		} else {
			steps = append(steps, stepStyle.Render(strconv.Itoa(i+1)))
		}
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		steps...,
	)
}
