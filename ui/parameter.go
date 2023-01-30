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
	paramTrackTitleStyle = lipgloss.NewStyle().
				Padding(1, 1, 1, 2).
				MarginRight(2).
				BorderStyle(lipgloss.HiddenBorder())

	paramStepTitleStyle = paramTrackTitleStyle.Copy().
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(primaryColor)

	paramStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2).
			Bold(true).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(secondaryColor)

	selectedParamStyle = paramStyle.Copy().
				BorderStyle(lipgloss.ThickBorder()).
				Foreground(primaryTextColor)
)

type parameters struct {
	track []parameter[sequencer.Track]
	step  []parameter[sequencer.Step]
}

type parameter[t sequencer.Parametrable] struct {
	value  func(item t) int
	string func(item t) string
	set    func(item t, value int, add int)
}

func (m *mainModel) initParameters() {
	m.parameters.track = []parameter[sequencer.Track]{
		{
			value: func(item sequencer.Track) int {
				// TODO: make chords actual chords
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.ChordString()),
					"",
					"note",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetChord([]uint8{
					uint8(value + add),
				})
			},
		},
		{
			value: func(item sequencer.Track) int {
				return item.Length()
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.LengthString()),
					"",
					"length",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				setLengthParam(item, value, add)
			},
		},
		{
			value: func(item sequencer.Track) int {
				return int(item.Velocity())
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.VelocityString()),
					"",
					"velocity",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetVelocity(uint8(value + add))
			},
		},
		{
			value: func(item sequencer.Track) int {
				return item.Probability()
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.ProbabilityString()),
					"",
					"probability",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetProbability(value + add)
			},
		},
		{
			value: func(item sequencer.Track) int {
				return item.Device()
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					"",
					item.DeviceString(),
					"",
					"",
					"device",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetDevice(value + add)
			},
		},
		{
			value: func(item sequencer.Track) int {
				return int(item.Channel())
			},
			string: func(item sequencer.Track) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.ChannelString()),
					"",
					"midi channel",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetChannel(uint8(value + add))
			},
		},
	}

	m.parameters.step = []parameter[sequencer.Step]{
		{
			value: func(item sequencer.Step) int {
				return int(item.Chord()[0])
			},
			string: func(item sequencer.Step) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.ChordString()),
					"",
					"note",
				)
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetChord([]uint8{
					uint8(value + add),
				})
			},
		},
		{
			value: func(item sequencer.Step) int {
				return item.Length()
			},
			string: func(item sequencer.Step) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.LengthString()),
					"",
					"length",
				)
			},
			set: func(item sequencer.Step, value int, add int) {
				setLengthParam(item, value, add)
			},
		},
		{
			value: func(item sequencer.Step) int {
				return int(item.Velocity())
			},
			string: func(item sequencer.Step) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.VelocityString()),
					"",
					"velocity",
				)
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetVelocity(uint8(value + add))
			},
		},
		{
			value: func(item sequencer.Step) int {
				return item.Probability()
			},
			string: func(item sequencer.Step) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.ProbabilityString()),
					"",
					"probability",
				)
			},
			set: func(item sequencer.Step, value int, add int) {
				item.SetProbability(value + add)
			},
		},
		{
			value: func(item sequencer.Step) int {
				return item.Offset()
			},
			string: func(item sequencer.Step) string {
				return lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(item.OffsetString()),
					"",
					"offset",
				)
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

func (m mainModel) renderParams() string {
	var params []string
	if m.mode == stepMode {
		params = append(params, paramStepTitleStyle.Render(toASCIIFont(fmt.Sprintf("S%d", m.activeStep+1))))
		if m.getActiveStep().IsActive() {
			for i, p := range m.parameters.step {
				var style lipgloss.Style
				if m.activeParam == i {
					style = selectedParamStyle
				} else {
					style = paramStyle
				}
				params = append(
					params,
					style.Render(p.string(m.getActiveStep())),
				)
			}
		}
	} else {
		params = append(params, paramTrackTitleStyle.Render(toASCIIFont(fmt.Sprintf("T%d", m.activeTrack+1))))
		for i, p := range m.parameters.track {
			var style lipgloss.Style
			if m.activeParam == i {
				style = selectedParamStyle
			} else {
				style = paramStyle
			}
			params = append(
				params,
				style.Render(p.string(m.getActiveTrack())),
			)
		}
	}
	return lipgloss.NewStyle().MarginTop(1).Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		params...,
	))
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
