package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	stepWidth  = 15
	stepHeight = stepWidth / 2

	stepCurrentColor        = lipgloss.Color("15")
	stepActiveColor         = lipgloss.Color("250")
	stepInactiveColor       = lipgloss.Color("240")
	stepTextBackgroundColor = lipgloss.Color("240")
	stepTextColor           = lipgloss.Color("232")

	stepStyle = lipgloss.NewStyle().Margin(1, 2, 0, 0)
	textStyle = lipgloss.NewStyle().
			Foreground(stepTextColor).
			Padding(1, 2).
			Bold(true)
)

func (m mainModel) renderStep(step *sequencer.Step) string {
	content := m.renderStepContent(step)
	width, height := m.stepSize()

	if step.IsCurrentStep() {
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

func (m mainModel) renderStepContent(step *sequencer.Step) string {
	return strconv.Itoa(step.Number() + 1)
}
