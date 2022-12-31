package ui

import (
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
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

func (u UI) renderStep(step *sequencer.Step) string {
	content := u.renderStepContent(step)
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

func (u UI) renderStepContent(step *sequencer.Step) string {
	return strconv.Itoa(step.Number() + 1)
}
