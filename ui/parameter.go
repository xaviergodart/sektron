package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/lipgloss"
)

const (
	paramsPerLine = 18
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
	set    func(item t, value int)
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
			set: func(item sequencer.Track, value int) {
				item.SetDevice(value)
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
			set: func(item sequencer.Track, value int) {
				item.SetChannel(uint8(value))
			},
		},
		{
			name: "note",
			value: func(item sequencer.Track) int {
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Track) string {
				return item.ChordString()
			},
			set: func(item sequencer.Track, value int) {
				item.SetChord([]uint8{
					uint8(value),
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
			set: func(item sequencer.Track, value int) {
				// TODO: update knob course
				item.SetLength(value)
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
			set: func(item sequencer.Track, value int) {
				item.SetVelocity(uint8(value))
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
			set: func(item sequencer.Track, value int) {
				item.SetProbability(value)
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
			set: func(item sequencer.Step, value int) {
				item.SetChord([]uint8{
					uint8(value),
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
			set: func(item sequencer.Step, value int) {
				item.SetLength(value)
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
			set: func(item sequencer.Step, value int) {
				item.SetVelocity(uint8(value))
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
			set: func(item sequencer.Step, value int) {
				item.SetProbability(value)
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
			set: func(item sequencer.Step, value int) {
				item.SetOffset(value)
			},
		},
	}
}

func (p *parameter[t]) update(item t, add int) {
	newValue := p.value(item) + add
	p.set(item, newValue)
}

func (p parameter[t]) render(item t) string {
	return fmt.Sprintf("%s: %s", p.name, p.string(item))
}

func (m mainModel) renderParams() string {
	params := make([]string, len(m.parameters.track))
	width, height := m.paramSize()
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
	return width, height
}
