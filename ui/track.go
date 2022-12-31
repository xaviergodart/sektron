package ui

import (
	"fmt"
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	trackTitle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Margin(1, 1).
			Bold(true)

	stepsPerPage = 16
	stepsPerLine = 8

	stepWith     = 14
	stepHeight   = 6
	primaryColor = lipgloss.Color("201")

	stepStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Width(stepWith).
			Height(stepHeight)
	stepStyleCurrent = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder()).
				Bold(true).
				Inherit(stepStyle)
	stepStyleActive = lipgloss.NewStyle().
			BorderForeground(primaryColor).
			Foreground(primaryColor).
			Inherit(stepStyle)
	stepStyleActiveCurrent = lipgloss.NewStyle().
				Inherit(stepStyleCurrent).
				Inherit(stepStyleActive)
)

func (u UI) viewTrack(track *sequencer.Track) string {
	pages := make([][]string, len(track.Steps()))

	for i, s := range track.Steps() {
		page := i / stepsPerPage
		var step string
		if !track.IsActive() {
			step = stepStyle.Render("")
			pages[page] = append(pages[page], step)
			continue
		}
		if i == track.CurrentStep() {
			if s.IsActive() {
				step = stepStyleActiveCurrent.Render(strconv.Itoa(i + 1))
			} else {
				step = stepStyleCurrent.Render(strconv.Itoa(i + 1))
			}
		} else {
			if s.IsActive() {
				step = stepStyleActive.Render(strconv.Itoa(i + 1))
			} else {
				step = stepStyle.Render(strconv.Itoa(i + 1))
			}
		}
		pages[page] = append(pages[page], step)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		u.viewTrackTitle(track),
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

func (u UI) viewTrackTitle(track *sequencer.Track) string {
	var title string
	value := fmt.Sprintf("TRACK: %d", u.activeTrack+1)
	if track.IsActive() {
		title = trackTitle.Strikethrough(false).Render(value)
	} else {
		title = trackTitle.Strikethrough(true).Render(value)
	}
	return title
}
