package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	mainView = lipgloss.NewStyle().
		Align(lipgloss.Center)
)

func (m mainModel) renderSequencer() string {
	track := m.seq.Tracks()[m.activeTrack]

	return mainView.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.renderTrack(track),
			),
		)
}
