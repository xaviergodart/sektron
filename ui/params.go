package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

const (
	paramsPerLine = 3
	paramWidth    = stepWidth * 2
	paramHeight   = stepWidth / 2
)

var (
	paramStyle = lipgloss.NewStyle().
		Margin(1, 2, 0, 0).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))
	/*textStyle  = lipgloss.NewStyle().
	Foreground(secondaryTextColor).
	Padding(1, 2).
	Bold(true)*/
)

type parameter struct {
	name        string
	value       func(track int, step int, value int) int
	updateTrack func(track int, value int)
	updateStep  func(track int, step int, value int)
}

func parameters(seq sequencer.Sequencer) []parameter {
	return []parameter{
		parameter{
			name: "chord",
			value: func() {
				return 0
			},
			updateTrack: func(track int, value int) {
				seq.Tracks()[track].SetChord([]uint8{
					uint8(value),
				})
			},
			updateStep: func(track int, step int, value int) {
				seq.Tracks()[track].SetChord([]uint8{
					uint8(value),
				})
			},
		},
	}
}

func (p *parameter) update(track int, step *int, add int) {
	if step == nil {
		p.updateTrack(track, p.value()+add)
	}
	p.updateStep(track, *step, p.value()+add)
}

func (p parameter) render() string {
	return fmt.Sprintf("%s: %d", p.name, p.value)
}

func (m mainModel) renderParams() string {
	params := make([]string, len(m.parameters))
	width, height := m.paramSize()
	for _, p := range m.parameters {
		params = append(
			params,
			paramStyle.Render(
				lipgloss.Place(
					width,
					height,
					lipgloss.Center,
					lipgloss.Center,
					p.render(),
				),
			),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		params...,
	)
}

func (m mainModel) paramSize() (int, int) {
	width := m.width / paramsPerLine
	height := width / 4
	if width < paramWidth || height < paramHeight {
		return paramWidth, paramHeight
	}
	return width, height
}
