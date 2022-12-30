package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	trackTitle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Align(lipgloss.Center).
		Padding(1, 1).
		Height(3).
		Margin(0, 1).
		Bold(true)
)

func (u UI) ViewSequencer() string {
	track := u.seq.Tracks()[u.activeTrack]
	return lipgloss.JoinVertical(
		lipgloss.Left,
		trackTitle.Render(fmt.Sprintf("TRACK: %d", u.activeTrack+1)),
		u.ViewTrack(track),
	)
}
