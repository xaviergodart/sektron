package ui

import (
	"fmt"
	"sektron/sequencer"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	stepWidth  = 15
	stepHeight = stepWidth / 2

	maxVelocityValue = 127
)

var (
	stepStyle = lipgloss.NewStyle().
			Margin(1, 2, 0, 0)
	stepSelectedStyle = lipgloss.NewStyle().
				Margin(1, 2, 0, 0).
				Foreground(primaryColor)
	stepActiveStyle = lipgloss.NewStyle().
			Margin(1, 0, 0, 0)
	stepVelocityStyle = lipgloss.NewStyle().
				Margin(1, 1, 0, 0).
				Foreground(secondaryColor)
	textStyle = lipgloss.NewStyle().
			Foreground(secondaryTextColor).
			Padding(1, 1, 1, 2).
			Bold(true)
)

func (m mainModel) renderStep(step sequencer.Step) string {
	content := m.renderStepContent(step)
	var stepStr string
	width, height := m.stepSize()

	var stepCurrentColor, stepActiveColor, stepInactiveColor lipgloss.Color
	if m.mode == stepMode && m.activeStep == step.Position() {
		stepCurrentColor = currentColor
		stepActiveColor = primaryColor
		stepInactiveColor = primaryColor
	} else if step.Track().IsActive() {
		stepCurrentColor = currentColor
		stepActiveColor = activeColor
		stepInactiveColor = inactiveColor
	} else {
		stepCurrentColor = currentDimmedColor
		stepActiveColor = activeDimmedColor
		stepInactiveColor = inactiveDimmedColor
	}

	if m.seq.IsPlaying() && step.IsCurrentStep() && step.IsActive() {
		stepStr = stepActiveStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(stepCurrentColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepCurrentColor),
		))
	} else if m.seq.IsPlaying() && step.IsCurrentStep() {
		stepStr = stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(stepCurrentColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepCurrentColor),
		))
	} else if step.IsActive() {
		stepStr = stepActiveStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(stepActiveColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepActiveColor),
		))

	} else {
		stepStr = stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(stepInactiveColor).Render(content),
			lipgloss.WithWhitespaceBackground(stepInactiveColor),
		))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		stepStr,
		m.renderVelocity(step, lipgloss.Height(stepStr)),
	)
}

func (m mainModel) stepSize() (int, int) {
	width := m.width/stepsPerLine - 2
	height := width/2 - 1
	if width < stepWidth || height < stepHeight {
		return stepWidth, stepHeight
	}
	return width, height
}

func (m mainModel) renderStepContent(step sequencer.Step) string {
	if !step.IsActive() {
		return strconv.Itoa(step.Position() + 1)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%d", step.Position()+1),
		toASCIIFont(step.ChordString()),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			step.LengthString(),
			lipgloss.NewStyle().
				MarginLeft(1).
				Render(step.ProbabilityString()),
			lipgloss.NewStyle().
				MarginLeft(1).
				Render(step.OffsetString()),
		),
	)
}

func (m mainModel) renderVelocity(step sequencer.Step, height int) string {
	if !step.IsActive() {
		return ""
	}
	velocityIndicator := []string{}
	velocityValue := int(maxVelocityValue-step.Velocity()) * height / maxVelocityValue
	for i := 1; i < height; i++ {
		if velocityValue < i {
			velocityIndicator = append(velocityIndicator, "█")
		} else {
			velocityIndicator = append(velocityIndicator, " ")
		}
	}
	return stepVelocityStyle.Render(lipgloss.JoinVertical(lipgloss.Left, velocityIndicator...))
}
