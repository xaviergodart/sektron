package ui

import (
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

var (
	stepsPerPage = 16
	stepsPerLine = 8
)

func (m mainModel) renderTrack(track *sequencer.Track) string {
	pages := make([][]string, len(track.Steps()))

	for i, step := range track.Steps() {
		page := i / stepsPerPage
		pages[page] = append(pages[page], m.renderStep(step))
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			pages[m.activeTrackPage][:stepsPerLine]...,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			pages[m.activeTrackPage][len(pages[m.activeTrackPage])-stepsPerLine:]...,
		),
	)
}
