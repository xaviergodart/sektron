package sequencer

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

var (
	stepStyle = lipgloss.NewStyle().
			Width(6).
			Height(3).
			Margin(1).
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))
	stepCurrentStyle = lipgloss.NewStyle().
				Width(6).
				Height(3).
				Margin(1).
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFFFFF"))
)

func (t track) View() string {
	ui := fmt.Sprintf("pulse: %d\n", t.pulse)
	var steps []string
	for i := range t.steps {
		if i == t.activeStep() {
			steps = append(steps, stepCurrentStyle.Render(strconv.Itoa(i+1)))
		} else {
			steps = append(steps, stepStyle.Render(strconv.Itoa(i+1)))
		}
	}
	return ui + lipgloss.JoinHorizontal(
		lipgloss.Left,
		steps...,
	)
}
