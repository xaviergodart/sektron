package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (u UI) ViewSequencer() string {
	track := u.seq.Tracks()[u.activeTrack]

	return lipgloss.JoinVertical(
		lipgloss.Left,
		u.viewTrack(track),
	)
}
