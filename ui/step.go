package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	stepWidth  = 13
	stepHeight = stepWidth / 2

	stepCurrentColor  = lipgloss.Color("15")
	stepActiveColor   = lipgloss.Color("250")
	stepInactiveColor = lipgloss.Color("240")

	stepCurrentInactiveTrackColor  = lipgloss.Color("240")
	stepActiveInactiveTrackColor   = lipgloss.Color("239")
	stepInactiveInactiveTrackColor = lipgloss.Color("234")

	stepTextBackgroundColor = lipgloss.Color("240")
	stepTextColor           = lipgloss.Color("232")
)

var (
	stepStyle = lipgloss.NewStyle().Margin(1, 2, 0, 0)
	textStyle = lipgloss.NewStyle().
			Foreground(stepTextColor).
			Padding(1, 2).
			Bold(true)
)

func (m mainModel) renderStep(step *sequencer.Step) string {
	content := m.renderStepContent(step)
	width, height := m.stepSize()

	var currentColor, activeColor, inactiveColor lipgloss.Color
	if step.Track().IsActive() {
		currentColor = stepCurrentColor
		activeColor = stepActiveColor
		inactiveColor = stepInactiveColor
	} else {
		currentColor = stepCurrentInactiveTrackColor
		activeColor = stepActiveInactiveTrackColor
		inactiveColor = stepInactiveInactiveTrackColor
	}

	if m.seq.IsPlaying() && step.IsCurrentStep() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			textStyle.Background(currentColor).Render(content),
			lipgloss.WithWhitespaceBackground(currentColor),
		))
	}
	if step.IsActive() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			textStyle.Background(activeColor).Render(content),
			lipgloss.WithWhitespaceBackground(activeColor),
		))
	}

	return stepStyle.Render(lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		textStyle.Background(inactiveColor).Render(content),
		lipgloss.WithWhitespaceBackground(inactiveColor),
	))
}

func (m mainModel) stepSize() (int, int) {
	width := m.width/stepsPerLine - 2
	height := width/2 - 1
	if width < stepWidth || height < stepHeight {
		return stepWidth, stepHeight
	}
	return width, height
}

func (m mainModel) renderStepContent(step *sequencer.Step) string {
	return strconv.Itoa(step.Number() + 1)
}
