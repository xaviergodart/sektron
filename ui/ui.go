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

	// recMode allows the user to activate/deactivate steps using the step keys.
	recMode
)

const (
	// We don't need to refresh the ui as often as the sequencer.
	// It saves some cpu. Right now we run it at 30 fps.
	refreshFrequency = 33 * time.Millisecond
)

type mainModel struct {
	seq             sequencer.Sequencer
	parameters      parameters
	keymap          keyMap
	width           int
	height          int
	mode            mode
	activeTrack     int
	activeTrackPage int
	activeStep      int
	activeParam     int
	help            help.Model
}

// New creates a new mainModel that hols the ui state. It takes a new sequencer.
// Check teh sequencer package.
func New(seq sequencer.Sequencer) mainModel {
	model := mainModel{
		seq:             seq,
		keymap:          DefaultKeyMap(),
		activeTrack:     0,
		activeTrackPage: 0,
		activeStep:      0,
		activeParam:     0,
		help:            help.New(),
	}
	model.initParameters()
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
		return m, nil

	case tickMsg:
		return m, tick()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.TogglePlay):
			m.seq.TogglePlay()
			return m, nil

		case key.Matches(msg, m.keymap.Mode):
			m.activeParam = 0
			m.activeStep = 0
			if m.mode == trackMode {
				m.mode = recMode
			} else {
				m.mode = trackMode
			}
			return m, nil

		case key.Matches(msg, m.keymap.Add):
			m.addPress(msg)
			return m, nil

		case key.Matches(msg, m.keymap.Remove):
			m.removePress(msg)
			return m, nil

		case key.Matches(msg, m.keymap.Select):
			m.selectPress(msg)
			return m, nil

		case key.Matches(msg, m.keymap.Toggle):
			number := m.keymap.ToggleIndex[msg.String()]
			if m.mode == trackMode {
				m.seq.ToggleTrack(number)
			} else if m.mode == recMode {
				m.seq.ToggleStep(m.activeTrack, number+(m.activeTrackPage*stepsPerPage))
			}
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

		case key.Matches(msg, m.keymap.Params):
			m.activeParam = m.keymap.ParamsIndex[msg.String()]
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueUp):
			// TODO: ugly move that to specific function
			if m.mode == recMode {
				m.parameters.step[m.activeParam].update(m.getActiveStep(), 1)
			} else {
				m.parameters.track[m.activeParam].update(m.getActiveTrack(), 1)
			}
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueDown):
			if m.mode == recMode {
				m.parameters.step[m.activeParam].update(m.getActiveStep(), -1)
			} else {
				m.parameters.track[m.activeParam].update(m.getActiveTrack(), -1)
			}
			return m, nil

		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

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
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTransport(),
		m.renderSequencer(),
		m.renderParams(),
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

func (m *mainModel) selectPress(msg tea.KeyMsg) {
	number := m.keymap.SelectIndex[msg.String()]
	switch m.mode {
	case trackMode:
		if number >= len(m.seq.Tracks()) {
			return
		}
		m.activeTrack = number
		m.activeTrackPage = 0
		m.activeStep = 0
		m.activeParam = 0
	case recMode:
		m.activeStep = number + (m.activeTrackPage * stepsPerPage)
	}
}

func (m *mainModel) addPress(msg tea.KeyMsg) {
	switch m.mode {
	case trackMode:
		m.seq.AddTrack()
	case recMode:
		m.seq.AddStep(m.activeTrack)
	}
}

func (m *mainModel) removePress(msg tea.KeyMsg) {
	switch m.mode {
	case trackMode:
		if m.activeTrack > 0 && m.activeTrack == len(m.seq.Tracks())-1 {
			m.activeTrack--
		}
		m.seq.RemoveTrack()
	case recMode:
		remainingStepsInPage := (len(m.getActiveTrack().Steps()) - 1) % stepsPerPage
		if m.activeTrackPage > 0 && remainingStepsInPage == 0 {
			m.activeTrackPage--
		}
		m.seq.RemoveStep(m.activeTrack)
	}
}

func (m *mainModel) getActiveTrack() sequencer.Track {
	return m.seq.Tracks()[m.activeTrack]
}

func (m *mainModel) getActiveStep() sequencer.Step {
	return m.seq.Tracks()[m.activeTrack].Steps()[m.activeStep]
}
