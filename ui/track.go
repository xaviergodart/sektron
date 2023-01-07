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
	pageNb := len(m.seq.Tracks()[m.activeTrack].Steps())/stepsPerPage + 1
	pages := make([][]string, pageNb)

	for i := range pages {
		pages[i] = make([]string, stepsPerPage)
	}

	for i, step := range track.Steps() {
		page := i / stepsPerPage
		pages[page][i%stepsPerPage] = m.renderStep(step)
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
