package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m mainModel) renderSequencer() string {
	track := m.seq.Tracks()[m.activeTrack]

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTrack(track),
	)
}
