package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m mainModel) renderSequencer() string {
	track := m.seq.Tracks()[m.activeTrack]

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderTrack(track),
	)
}
