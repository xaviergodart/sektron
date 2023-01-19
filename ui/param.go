package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

const (
	paramsPerLine = 6
	paramWidth    = stepWidth * 2
	paramHeight   = paramWidth / 2
)

var (
	paramStyle = lipgloss.NewStyle().
			Margin(1, 2, 0, 0).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(secondaryColor)

	selectedParamStyle = paramStyle.Copy().
				BorderStyle(lipgloss.ThickBorder())
)

type parameter struct {
	name        string
	min         int
	max         int
	value       func(item sequencer.Parametrable) int
	string      func(item sequencer.Parametrable) string
	updateValue func(item sequencer.Parametrable, value int)
}

func (m *mainModel) initParameters() {
	m.parameters = []parameter{
		{
			name: "note",
			min:  21,
			max:  108,
			value: func(item sequencer.Parametrable) int {
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Parametrable) string {
				return sequencer.ChordString(item.Chord())
			},
			updateValue: func(item sequencer.Parametrable, value int) {
				item.SetChord([]uint8{
					uint8(value),
				})
			},
		},
		{
			name: "length",
			min:  1,
			max:  128 * 6,
			value: func(item sequencer.Parametrable) int {
				return item.Length()
			},
			string: func(item sequencer.Parametrable) string {
				return sequencer.LengthString(item.Length())
			},
			updateValue: func(item sequencer.Parametrable, value int) {
				item.SetLength(value)
			},
		},
		/*		{
					name: "velocity",
					min:  1,
					max:  127,
					value: func(track int) int {
						return int(m.seq.Tracks()[track].Velocity())
					},
					string: func(track int) string {
						return fmt.Sprintf("%d", m.seq.Tracks()[track].Velocity())
					},
					updateTrack: func(track int, value int) {
						m.seq.Tracks()[track].SetVelocity(uint8(value))
					},
					updateStep: func(track int, step int, value int) {
						m.seq.Tracks()[track].Steps()[step].SetVelocity(uint8(value))
					},
				},
				{
					name: "probability",
					min:  1,
					max:  100,
					value: func(track int) int {
						return int(m.seq.Tracks()[track].Probability())
					},
					string: func(track int) string {
						return fmt.Sprintf("%d", m.seq.Tracks()[track].Probability())
					},
					updateTrack: func(track int, value int) {
						m.seq.Tracks()[track].SetProbability(value)
					},
					updateStep: func(track int, step int, value int) {
						m.seq.Tracks()[track].Steps()[step].SetProbability(value)
					},
				},*/
	}
}

func (p *parameter) update(item sequencer.Parametrable, add int) {
	newValue := p.value(item) + add
	if newValue < p.min || newValue > p.max {
		return
	}
	p.updateValue(item, newValue)
}

func (p parameter) render(item sequencer.Parametrable) string {
	return fmt.Sprintf("%s: %s", p.name, p.string(item))
}

func (m mainModel) renderParams() string {
	params := make([]string, len(m.parameters))
	width, height := m.paramSize()
	for i, p := range m.parameters {
		var style lipgloss.Style
		if m.activeParam == i {
			style = selectedParamStyle
		} else {
			style = paramStyle
		}
		params = append(
			params,
			style.Render(
				lipgloss.Place(
					width,
					height,
					lipgloss.Center,
					lipgloss.Center,
					p.render(m.seq.Tracks()[m.activeTrack]),
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
	height := width / 6
	if width < paramWidth || height < paramHeight {
		return paramWidth, paramHeight
	}
	return width, height
}
