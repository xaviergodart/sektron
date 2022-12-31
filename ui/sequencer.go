package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (u UI) renderSequencer() string {
	track := u.seq.Tracks()[u.activeTrack]

	return lipgloss.JoinVertical(
		lipgloss.Left,
		u.renderTrack(track),
	)
}
