// Packge ui handles the user interface for viewing and interacting with the
// sequencer package.
//
// It uses the bubbletea tui framework and lipgloss to make things look good.
package ui

import (
	"sektron/sequencer"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tickMsg is a message that triggers ui rrefresh
type tickMsg time.Time

// mode represents the different modes of the ui
type mode uint8

const (
	// trackMode allows the user to select the tracks using the step keys.
	trackMode mode = iota

	// stepMode allows the user to activate/deactivate steps using the step keys.
	stepMode

	paramSelectMode
)

const (
	// We don't need to refresh the ui as often as the sequencer.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond

	stepModeTimeout = 90
)

type mainModel struct {
	seq             sequencer.Sequencer
	parameters      parameters
	paramMidiTable  table.Model
	keymap          keyMap
	width           int
	height          int
	mode            mode
	activeTrack     int
	activeTrackPage int
	activeStep      int
	activeParams    []struct{ track, step int }
	stepModeTimer   int
	help            help.Model
}

// New creates a new mainModel that hols the ui state. It takes a new sequencer.
// Check teh sequencer package.
func New(seq sequencer.Sequencer) mainModel {
	model := mainModel{
		seq:          seq,
		keymap:       DefaultKeyMap(),
		activeParams: make([]struct{ track, step int }, 10),
		help:         help.New(),
	}
	rows := []table.Row{}
	for _, c := range seq.Tracks()[0].Controls() {
		rows = append(rows, table.Row{c.Name()})
	}
	model.initParameters()
	model.paramMidiTable = table.New(
		table.WithColumns([]table.Column{
			{Title: "midi message", Width: 40},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithKeyMap(table.DefaultKeyMap()),
	)
	return model
}

func tick() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.paramMidiTable.SetWidth(msg.Width)
		return m, nil

	case tickMsg:
		if m.mode == stepMode {
			m.stepModeTimer++
		}
		if m.stepModeTimer > stepModeTimeout {
			m.stepModeTimer = 0
			m.mode = trackMode
		}
		return m, tick()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.TogglePlay):
			m.seq.TogglePlay()
			return m, nil

		case key.Matches(msg, m.keymap.Mode):
			m.activeStep = 0
			if m.mode == trackMode {
				m.mode = stepMode
			} else {
				m.mode = trackMode
			}
			return m, nil

		case key.Matches(msg, m.keymap.Validate):
			if m.mode == paramSelectMode {
				m.getActiveTrack().AddControl(m.paramMidiTable.Cursor())
				m.mode = trackMode
			}
			return m, nil

		case key.Matches(msg, m.keymap.AddTrack):
			m.seq.AddTrack()
			return m, nil

		case key.Matches(msg, m.keymap.RemoveTrack):
			if m.activeTrack > 0 && m.activeTrack == len(m.seq.Tracks())-1 {
				m.activeTrack--
			}
			m.seq.RemoveTrack()
			return m, nil

		case key.Matches(msg, m.keymap.AddStep):
			m.seq.AddStep(m.activeTrack)
			return m, nil

		case key.Matches(msg, m.keymap.RemoveStep):
			remainingStepsInPage := (len(m.getActiveTrack().Steps()) - 1) % stepsPerPage
			if m.activeTrackPage > 0 && remainingStepsInPage == 0 {
				m.activeTrackPage--
			}
			m.activeStep = 0
			m.seq.RemoveStep(m.activeTrack)
			return m, nil

		case key.Matches(msg, m.keymap.StepSelect):
			number := m.keymap.StepSelectIndex[msg.String()]
			if number >= len(m.getActiveTrack().Steps())-(m.activeTrackPage*stepsPerPage) {
				return m, nil
			}
			m.activeStep = number + (m.activeTrackPage * stepsPerPage)
			m.mode = stepMode
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.StepToggle):
			number := m.keymap.StepToggleIndex[msg.String()]
			if number >= len(m.getActiveTrack().Steps()) {
				return m, nil
			}
			m.activeStep = number + (m.activeTrackPage * stepsPerPage)
			m.seq.ToggleStep(m.activeTrack, m.activeStep)
			m.mode = stepMode
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.TrackSelect):
			number := m.keymap.TrackSelectIndex[msg.String()]
			if number >= len(m.seq.Tracks()) {
				return m, nil
			}
			m.activeTrack = number
			m.activeTrackPage = 0
			m.activeStep = 0
			return m, nil

		case key.Matches(msg, m.keymap.TrackToggle):
			number := m.keymap.TrackToggleIndex[msg.String()]
			m.seq.ToggleTrack(number)
			return m, nil

		case key.Matches(msg, m.keymap.TrackPageUp):
			pageNb := m.trackPagesNb()
			m.activeTrackPage = (m.activeTrackPage + 1) % pageNb
			return m, nil

		case key.Matches(msg, m.keymap.TrackPageDown):
			pageNb := m.trackPagesNb()
			if m.activeTrackPage-1 < 0 {
				m.activeTrackPage = pageNb - 1
			} else {
				m.activeTrackPage = m.activeTrackPage - 1
			}
			return m, nil

		case key.Matches(msg, m.keymap.TempoUp):
			m.seq.SetTempo(m.seq.Tempo() + 1)
			return m, nil

		case key.Matches(msg, m.keymap.TempoDown):
			m.seq.SetTempo(m.seq.Tempo() - 1)
			return m, nil

		case key.Matches(msg, m.keymap.TempoFineUp):
			m.seq.SetTempo(m.seq.Tempo() + 0.1)
			return m, nil

		case key.Matches(msg, m.keymap.TempoFineDown):
			m.seq.SetTempo(m.seq.Tempo() - 0.1)
			return m, nil

		case key.Matches(msg, m.keymap.AddParam):
			m.mode = paramSelectMode
			return m, nil

		case key.Matches(msg, m.keymap.RemoveParam):
			// TODO: Not working
			m.getActiveTrack().RemoveControl(m.getActiveParam())
			return m, nil

		case key.Matches(msg, m.keymap.ParamSelectLeft):
			// TODO: Not working with added parameters
			m.setActiveParam(m.getActiveParam() - 1)
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamSelectRight):
			m.setActiveParam(m.getActiveParam() + 1)
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueUp):
			if m.mode == stepMode && m.getActiveStep().IsActive() {
				m.parameters.step[m.getActiveParam()].increase(m.getActiveStep())
			} else if m.mode == trackMode {
				m.parameters.track[m.getActiveParam()].increase(m.getActiveTrack())
			} else if m.mode == paramSelectMode {
				var cmd tea.Cmd
				m.paramMidiTable, cmd = m.paramMidiTable.Update(msg)
				return m, cmd
			}
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueDown):
			if m.mode == stepMode && m.getActiveStep().IsActive() {
				m.parameters.step[m.getActiveParam()].decrease(m.getActiveStep())
			} else if m.mode == trackMode {
				m.parameters.track[m.getActiveParam()].decrease(m.getActiveTrack())
			} else if m.mode == paramSelectMode {
				var cmd tea.Cmd
				m.paramMidiTable, cmd = m.paramMidiTable.Update(msg)
				return m, cmd
			}
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.ClearScreen

		case key.Matches(msg, m.keymap.Quit):
			if m.seq.IsPlaying() {
				m.seq.TogglePlay()
			}
			m.seq.Reset()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	var params string
	if m.mode == paramSelectMode {
		params = m.paramMidiTable.View()
	} else {
		params = m.renderParams()
	}
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTransport(),
		m.renderSequencer(),
		params,
	)

	help := m.help.View(m.keymap)

	// Cleanup gibber
	cleanup := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - lipgloss.Height(mainView) - lipgloss.Height(help)).
		Render("")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		cleanup,
		help,
	)
}

func (m *mainModel) getActiveTrack() sequencer.Track {
	return m.seq.Tracks()[m.activeTrack]
}

func (m *mainModel) getActiveStep() sequencer.Step {
	return m.seq.Tracks()[m.activeTrack].Steps()[m.activeStep]
}

func (m mainModel) getActiveParam() int {
	if m.mode == stepMode {
		return m.activeParams[m.activeTrack].step
	}
	return m.activeParams[m.activeTrack].track
}

func (m *mainModel) setActiveParam(value int) {
	min := 0
	if m.mode == stepMode {
		max := m.stepParamCount()
		if value >= min && value < max {
			m.activeParams[m.activeTrack].step = value
		}
		return
	}
	max := m.trackParamCount()
	if value >= min && value < max {
		m.activeParams[m.activeTrack].track = value
	}
}
