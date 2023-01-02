package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	stepWidth  = 8
	stepHeight = stepWidth/2 - 1

	stepStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Width(stepWidth).
			Height(stepHeight)
	stepStyleCurrent = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(focusColor).
				Bold(true).
				Inherit(stepStyle)
	stepStyleActive = lipgloss.NewStyle().
			BorderForeground(primaryColor).
			Foreground(primaryColor).
			Inherit(stepStyle)
	stepStyleActiveCurrent = lipgloss.NewStyle().
				Inherit(stepStyleCurrent).
				UnsetBorderForeground().
				Inherit(stepStyleActive)
)

func (m mainModel) renderStep(step *sequencer.Step) string {
	content := m.renderStepContent(step)
	width, height := m.stepSize()
	return m.stepStyle(step).
		Width(width).
		Height(height).
		Render(content)
}

func (m mainModel) stepSize() (int, int) {
	width := m.width/stepsPerLine - 2
	height := width/2 - 1
	if width < stepWidth || height < stepHeight {
		return stepWidth, stepHeight
	}
	return width, height
}

func (m mainModel) stepStyle(step *sequencer.Step) lipgloss.Style {
	if !step.Track().IsActive() {
		return stepStyle
	}
	if step.IsCurrentStep() {
		if step.IsActive() {
			return stepStyleActiveCurrent
		}
		return stepStyleCurrent
	}
	if step.IsActive() {
		return stepStyleActive
	}
	return stepStyle
}

func (m mainModel) renderStepContent(step *sequencer.Step) string {
	return strconv.Itoa(step.Number() + 1)
}
