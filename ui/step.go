package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	stepWidth  = 13
	stepHeight = stepWidth / 2
)

var (
	stepStyle = lipgloss.NewStyle().Margin(1, 2, 0, 0)
	textStyle = lipgloss.NewStyle().
			Foreground(secondaryTextColor).
			Padding(1, 2).
			Bold(true)
)

func (m mainModel) renderStep(step sequencer.StepInterface) string {
	content := m.renderStepContent(step)
	width, height := m.stepSize()

	var stepCurrentColor, stepActiveColor, stepInactiveColor lipgloss.Color
	if step.Track().IsActive() {
		stepCurrentColor = currentColor
		stepActiveColor = activeColor
		stepInactiveColor = inactiveColor
	} else {
		stepCurrentColor = currentDimmedColor
		stepActiveColor = activeDimmedColor
		stepInactiveColor = inactiveDimmedColor
	}

	if m.seq.IsPlaying() && step.IsCurrentStep() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			textStyle.Background(stepCurrentColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepCurrentColor),
		))
	}
	if step.IsActive() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			textStyle.Background(stepActiveColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepActiveColor),
		))
	}

	return stepStyle.Render(lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		textStyle.Background(stepInactiveColor).Render(content),
		lipgloss.WithWhitespaceBackground(stepInactiveColor),
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

func (m mainModel) renderStepContent(step sequencer.StepInterface) string {
	return strconv.Itoa(step.Number() + 1)
}
