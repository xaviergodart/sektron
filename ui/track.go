package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

var (
	trackTitle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Margin(1, 1).
			Bold(true)

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
		m.renderTrackTitle(track),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			pages[m.activeTrackPage][:stepsPerLine]...,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			pages[m.activeTrackPage][len(pages[m.activeTrackPage])-stepsPerLine:]...,
		),
	)
}

func (m mainModel) renderTrackTitle(track *sequencer.Track) string {
	var title string
	value := fmt.Sprintf("TRACK: %d", m.activeTrack+1)
	if track.IsActive() {
		title = trackTitle.Strikethrough(false).Render(value)
	} else {
		title = trackTitle.Strikethrough(true).Render(value)
	}
	return title
}
