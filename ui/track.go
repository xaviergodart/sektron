package ui

import (
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

const (
	stepsPerPage = 16
	stepsPerLine = 8
)

func (m mainModel) renderTrack(track sequencer.Track) string {
	pages := make([][]string, m.trackPagesNb())

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

func (m mainModel) trackPagesNb() int {
	pageNb := len(m.getActiveTrack().Steps()) / stepsPerPage
	if len(m.getActiveTrack().Steps())%stepsPerPage > 0 {
		pageNb++
	}
	return pageNb
}
