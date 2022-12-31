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

func (u UI) renderTrack(track *sequencer.Track) string {
	pages := make([][]string, len(track.Steps()))

	for i, step := range track.Steps() {
		page := i / stepsPerPage
		pages[page] = append(pages[page], u.renderStep(step))
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		u.renderTrackTitle(track),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			pages[u.activeTrackPage][:stepsPerLine]...,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			pages[u.activeTrackPage][len(pages[u.activeTrackPage])-stepsPerLine:]...,
		),
	)
}

func (u UI) renderTrackTitle(track *sequencer.Track) string {
	var title string
	value := fmt.Sprintf("TRACK: %d", u.activeTrack+1)
	if track.IsActive() {
		title = trackTitle.Strikethrough(false).Render(value)
	} else {
		title = trackTitle.Strikethrough(true).Render(value)
	}
	return title
}
