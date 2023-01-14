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
)

type parameter struct {
	name        string
	min         int
	max         int
	value       func(track int) int
	updateTrack func(track int, value int)
	updateStep  func(track int, step int, value int)
}

func parameters(seq sequencer.Sequencer) []parameter {
	return []parameter{
		{
			name: "chord",
			min:  21,
			max:  108,
			value: func(track int) int {
				return int(seq.Tracks()[track].Chord()[0])
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
	newValue := p.value(track) + add
	if newValue < p.min || newValue > p.max {
		return
	}
	if step == nil {
		p.updateTrack(track, newValue)
	} else {
		p.updateStep(track, *step, newValue)
	}
}

func (p parameter) render(track int) string {
	return fmt.Sprintf("%s: %d", p.name, p.value(track))
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
					p.render(m.activeTrack),
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
