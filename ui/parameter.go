package ui

import (
	"fmt"
	"sektron/sequencer"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	carousel "github.com/xaviergodart/bubble-carousel"
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
				MarginTop(1).
				Bold(true).
				BorderStyle(lipgloss.HiddenBorder())

	paramStepTitleStyle = paramTrackTitleStyle.Copy().
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(primaryColor)

	paramStyles = carousel.Styles{
		Item: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2).
			BorderStyle(lipgloss.HiddenBorder()),
		Selected: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(secondaryColor),
	}
)

type parameters struct {
	track        []parameter[sequencer.Track]
	step         []parameter[sequencer.Step]
	index        map[int]int
	title        string
	content      string
	fixedParamNb int
}

func (p *parameters) getStepParam(nb int) *parameter[sequencer.Step] {
	return &p.step[p.index[nb]]
}

func (p *parameters) getTrackParam(nb int) *parameter[sequencer.Track] {
	return &p.track[p.index[nb]]
}

func (p parameters) getParamIndex(nb int) int {
	return p.index[nb]
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
	m.paramCarousel = carousel.New(
		carousel.WithFocused(true),
		carousel.WithStyles(paramStyles),
	)

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

	m.updateParams()
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
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.parameters.title,
		m.parameters.content,
	)
}

func (m *mainModel) updateParams() {
	switch m.mode {
	case stepMode:
		m.parameters.title = paramStepTitleStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				toASCIIFont(fmt.Sprintf("S%d", m.activeStep+1)),
				"",
				fmt.Sprintf("pattern %d", m.seq.ActivePattern()+1),
			),
		)
	case trackMode, paramSelectMode:
		m.parameters.title = paramTrackTitleStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				toASCIIFont(fmt.Sprintf("T%d", m.activeTrack+1)),
				"",
				fmt.Sprintf("pattern %d", m.seq.ActivePattern()+1),
			),
		)
	default:
		m.parameters.title = ""
	}

	var params []string
	m.parameters.index = map[int]int{}
	if m.mode == stepMode {
		if m.getActiveStep().IsActive() {
			for i, p := range m.parameters.step {
				if !p.active(m.getActiveStep()) {
					continue
				}
				m.parameters.index[len(params)] = i
				params = append(
					params,
					p.string(m.getActiveStep()),
				)
			}
		}
	} else if m.mode == trackMode {
		for i, p := range m.parameters.track {
			if !p.active(m.getActiveTrack()) {
				continue
			}
			m.parameters.index[len(params)] = i
			params = append(
				params,
				p.string(m.getActiveTrack()),
			)
		}
	} else if m.mode == paramSelectMode {
		scrollIndicator := []string{
			" ",
			"⏶",
			"⏷",
		}
		m.parameters.content = lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Top,
				scrollIndicator...,
			),
			m.paramMidiTable.View(),
		)
		return
	}
	m.paramCarousel.SetItems(params)
	m.paramCarousel.SetCursor(m.getActiveParam())
	m.parameters.content = lipgloss.NewStyle().
		MarginTop(1).
		Render(m.paramCarousel.View())
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
