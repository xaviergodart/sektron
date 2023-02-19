package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m mainModel) renderSequencer() string {
	if m.mode == patternSelectMode {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			m.renderPatterns(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderTrack(m.getActiveTrack()),
	)
}
