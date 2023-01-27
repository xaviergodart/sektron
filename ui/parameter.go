package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

const (
	paramsPerLine = 18
	pulsesPerStep = 6
	maxSteps      = 128
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

type parameters struct {
	track []parameter[sequencer.Track]
	step  []parameter[sequencer.Step]
}

type parameter[t sequencer.Parametrable] struct {
	name   string
	value  func(item t) int
	string func(item t) string
	set    func(item t, value int, add int)
}

func (m *mainModel) initParameters() {
	m.parameters.track = []parameter[sequencer.Track]{
		{
			name: "device",
			value: func(item sequencer.Track) int {
				return item.Device()
			},
			string: func(item sequencer.Track) string {
				return item.DeviceString()
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetDevice(value + add)
			},
		},
		{
			name: "channel",
			value: func(item sequencer.Track) int {
				return int(item.Channel())
			},
			string: func(item sequencer.Track) string {
				return item.ChannelString()
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetChannel(uint8(value + add))
			},
		},
		{
			name: "note",
			value: func(item sequencer.Track) int {
				// TODO: make chords actual chords
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Track) string {
				return item.ChordString()
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetChord([]uint8{
					uint8(value + add),
				})
			},
		},
		{
			name: "length",
			value: func(item sequencer.Track) int {
				return item.Length()
			},
			string: func(item sequencer.Track) string {
				return item.LengthString()
			},
			set: func(item sequencer.Track, value int, add int) {
				setLengthParam(item, value, add)
			},
		},
		{
			name: "velocity",
			value: func(item sequencer.Track) int {
				return int(item.Velocity())
			},
			string: func(item sequencer.Track) string {
				return item.VelocityString()
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetVelocity(uint8(value + add))
			},
		},
		{
			name: "probability",
			value: func(item sequencer.Track) int {
				return item.Probability()
			},
			string: func(item sequencer.Track) string {
				return item.ProbabilityString()
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetProbability(value + add)
			},
		},
	}

	m.parameters.step = []parameter[sequencer.Step]{
		{
			name: "note",
			value: func(item sequencer.Step) int {
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Step) string {
				return item.ChordString()
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetChord([]uint8{
					uint8(value + add),
				})
			},
		},
		{
			name: "length",
			value: func(item sequencer.Step) int {
				return item.Length()
			},
			string: func(item sequencer.Step) string {
				return item.LengthString()
			},
			set: func(item sequencer.Step, value int, add int) {
				setLengthParam(item, value, add)
			},
		},
		{
			name: "velocity",
			value: func(item sequencer.Step) int {
				return int(item.Velocity())
			},
			string: func(item sequencer.Step) string {
				return item.VelocityString()
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetVelocity(uint8(value + add))
			},
		},
		{
			name: "probability",
			value: func(item sequencer.Step) int {
				return item.Probability()
			},
			string: func(item sequencer.Step) string {
				return item.ProbabilityString()
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetProbability(value + add)
			},
		},
		{
			name: "offset",
			value: func(item sequencer.Step) int {
				return item.Offset()
			},
			string: func(item sequencer.Step) string {
				return item.OffsetString()
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetOffset(value + add)
			},
		},
	}
}

func (p *parameter[t]) increase(item t) {
	p.set(item, p.value(item), 1)
}

func (p *parameter[t]) decrease(item t) {
	p.set(item, p.value(item), -1)
}

func (p parameter[t]) render(item t) string {
	return fmt.Sprintf("%s: %s", p.name, p.string(item))
}

func (m mainModel) renderParams() string {
	var params []string
	width, height := m.paramSize()
	// TODO: ugly. Cleanup that...
	if m.mode == recMode {
		for i, p := range m.parameters.step {
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
						p.render(m.getActiveStep()),
					),
				),
			)
		}
	} else {
		for i, p := range m.parameters.track {
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
						p.render(m.getActiveTrack()),
					),
				),
			)
		}
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		params...,
	)
}

func (m mainModel) paramSize() (int, int) {
	width := m.width / paramsPerLine
	height := width / 6
	return width, height
}

func setLengthParam(item sequencer.Parametrable, value int, add int) {
	switch {
	case value < pulsesPerStep*4:
		item.SetLength(value + add)
	case value < pulsesPerStep*8:
		item.SetLength(value + 3*add)
	case value < pulsesPerStep*16:
		item.SetLength(value + 6*add)
	case value < pulsesPerStep*32:
		item.SetLength(value + 12*add)
	case value == pulsesPerStep*maxSteps+pulsesPerStep && add < 0:
		item.SetLength(pulsesPerStep * maxSteps)
	default:
		item.SetLength(value + 24*add)

	}
}
