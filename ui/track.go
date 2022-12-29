package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
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

func (u UI) ViewTrack(track *sequencer.Track) string {
	var steps []string
	for i := range track.Steps() {
		if i == track.ActiveStep() {
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
