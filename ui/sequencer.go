package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m mainModel) renderSequencer() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderTrack(m.getActiveTrack()),
	)
}
