package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const (
	paramsPerLine  = 18
	pulsesPerStep  = 6
	maxSteps       = 128
	midiParameters = 131
)

var (
	paramTrackTitleStyle = lipgloss.NewStyle().
				Padding(0, 1, 0, 2).
				MarginRight(2).
				Bold(true).
				BorderStyle(lipgloss.HiddenBorder())

	paramStepTitleStyle = paramTrackTitleStyle.Copy().
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(primaryColor)

	paramStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(secondaryColor)

	selectedParamStyle = paramStyle.Copy().
				BorderStyle(lipgloss.ThickBorder()).
				Foreground(primaryTextColor)
)

type parameters struct {
	track        []parameter[sequencer.Track]
	step         []parameter[sequencer.Step]
	fixedParamNb int
}

type parameter[t sequencer.Parametrable] struct {
	value  func(item t) int
	string func(item t) string
	set    func(item t, value int, add int)
	active func(item t) bool
}

func newMidiParameter[t sequencer.Parametrable](nb int) parameter[t] {
	return parameter[t]{
		value: func(item t) int {
			return int(item.Control(nb).Value())
		},
		string: func(item t) string {
			return lipgloss.JoinVertical(
				lipgloss.Center,
				toASCIIFont(item.Control(nb).String()),
				"",
				item.Control(nb).Name(),
			)
		},
		set: func(item t, value int, add int) {
			item.SetControl(nb, int16(value+add))
		},
		active: func(item t) bool {
			return item.IsActiveControl(nb)
		},
	}
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
			active: func(item sequencer.Track) bool {
				return true
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
			active: func(item sequencer.Track) bool {
				return true
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
			active: func(item sequencer.Track) bool {
				return true
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
			active: func(item sequencer.Track) bool {
				return true
			},
		},
		{
			value: func(item sequencer.Track) int {
				return item.Device()
			},
			string: func(item sequencer.Track) string {
				device := item.DeviceString()
				if len(device) > 40 {
					device = device[:37] + "..."
				}
				return lipgloss.JoinVertical(
					lipgloss.Center,
					"",
					wordwrap.String(device, 20),
					"",
					"device",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetDevice(value + add)
			},
			active: func(item sequencer.Track) bool {
				return true
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
					"channel",
				)
			},
			set: func(item sequencer.Track, value int, add int) {
				item.SetChannel(uint8(value + add))
			},
			active: func(item sequencer.Track) bool {
				return true
			},
		},
	}

	m.parameters.fixedParamNb = len(m.parameters.track)

	for i := 0; i <= midiParameters; i++ {
		m.parameters.track = append(m.parameters.track, newMidiParameter[sequencer.Track](i))
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
			active: func(item sequencer.Step) bool {
				return true
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
			active: func(item sequencer.Step) bool {
				return true
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
			active: func(item sequencer.Step) bool {
				return true
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
			active: func(item sequencer.Step) bool {
				return true
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
			active: func(item sequencer.Step) bool {
				return true
			},
		},
	}

	for i := 0; i <= midiParameters; i++ {
		m.parameters.step = append(m.parameters.step, newMidiParameter[sequencer.Step](i))
	}
}

func (m *mainModel) initMidiControls() {
	rows := []table.Row{}
	for _, c := range m.seq.Tracks()[0].Controls() {
		rows = append(rows, table.Row{c.Name()})
	}

	m.paramMidiTable = table.New(
		table.WithColumns([]table.Column{
			{Width: 40},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithKeyMap(table.DefaultKeyMap()),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(primaryTextColor).
		Background(secondaryColor).
		Bold(true)

	m.paramMidiTable.SetStyles(s)
}

func (p *parameter[t]) increase(item t) {
	p.set(item, p.value(item), 1)
}

func (p *parameter[t]) decrease(item t) {
	p.set(item, p.value(item), -1)
}

func (m mainModel) renderParams() string {
	var params []string
	var title string
	switch m.mode {
	case stepMode:
		title = paramStepTitleStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				toASCIIFont(fmt.Sprintf("S%d", m.activeStep+1)),
				"",
				fmt.Sprintf("pattern %d", m.seq.ActivePattern()+1),
			),
		)
	case trackMode, paramSelectMode:
		title = paramTrackTitleStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				toASCIIFont(fmt.Sprintf("T%d", m.activeTrack+1)),
				"",
				fmt.Sprintf("pattern %d", m.seq.ActivePattern()+1),
			),
		)
	default:
		title = ""
	}
	// TODO: render params on 2 lines
	if m.mode == stepMode {
		if m.getActiveStep().IsActive() {
			for i, p := range m.parameters.step {
				if !p.active(m.getActiveStep()) {
					continue
				}
				var style lipgloss.Style
				if m.getActiveParam() == i {
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
	} else if m.mode == trackMode {
		for i, p := range m.parameters.track {
			if !p.active(m.getActiveTrack()) {
				continue
			}
			var style lipgloss.Style
			if m.getActiveParam() == i {
				style = selectedParamStyle
			} else {
				style = paramStyle
			}
			params = append(
				params,
				style.Render(p.string(m.getActiveTrack())),
			)
		}
	} else if m.mode == paramSelectMode {
		scrollIndicator := []string{
			" ",
			"???",
			"???",
		}
		params = append(params,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Top,
					scrollIndicator...,
				),
				m.paramMidiTable.View(),
			),
		)
	}
	return lipgloss.NewStyle().
		MarginTop(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				append([]string{title}, params...)...,
			),
		)
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
