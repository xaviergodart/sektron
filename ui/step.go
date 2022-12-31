package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	stepWith     = 14
	stepHeight   = 6
	primaryColor = lipgloss.Color("207")
	focusColor   = lipgloss.Color("15")

	stepStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Width(stepWith).
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
	if !step.Track().IsActive() {
		return stepStyle.Render("")
	}
	if step.IsCurrentStep() {
		if step.IsActive() {
			return stepStyleActiveCurrent.Render(content)
		}
		return stepStyleCurrent.Render(content)
	}
	if step.IsActive() {
		return stepStyleActive.Render(content)
	}
	return stepStyle.Render(content)
}

func (m mainModel) renderStepContent(step *sequencer.Step) string {
	return strconv.Itoa(step.Number() + 1)
}
