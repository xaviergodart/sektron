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

	// stepMode allows the user to activate/deactivate steps using the step keys.
	stepMode
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
	keymap          keyMap
	width           int
	height          int
	mode            mode
	activeTrack     int
	activeTrackPage int
	activeStep      int
	activeParam     int
	stepModeTimer   int
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
		stepModeTimer:   0,
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
			m.activeParam = 0
			m.activeStep = 0
			if m.mode == trackMode {
				m.mode = stepMode
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

		case key.Matches(msg, m.keymap.StepSelect):
			number := m.keymap.StepSelectIndex[msg.String()]
			m.activeStep = number + (m.activeTrackPage * stepsPerPage)
			m.mode = stepMode
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.StepToggle):
			number := m.keymap.StepToggleIndex[msg.String()]
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
			m.activeParam = 0
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

		case key.Matches(msg, m.keymap.ParamSelectLeft):
			if m.activeParam > 0 {
				m.activeParam--
			}
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamSelectRight):
			max := 0
			if m.mode == trackMode {
				max = len(m.parameters.track) - 1
			} else {
				max = len(m.parameters.step) - 1
			}
			if m.activeParam < max {
				m.activeParam++
			}
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueUp):
			if m.mode == stepMode {
				m.parameters.step[m.activeParam].increase(m.getActiveStep())
			} else {
				m.parameters.track[m.activeParam].increase(m.getActiveTrack())
			}
			m.stepModeTimer = 0
			return m, nil

		case key.Matches(msg, m.keymap.ParamValueDown):
			if m.mode == stepMode {
				m.parameters.step[m.activeParam].decrease(m.getActiveStep())
			} else {
				m.parameters.track[m.activeParam].decrease(m.getActiveTrack())
			}
			m.stepModeTimer = 0
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

func (m *mainModel) addPress(msg tea.KeyMsg) {
	switch m.mode {
	case trackMode:
		m.seq.AddTrack()
	case stepMode:
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
	case stepMode:
		remainingStepsInPage := (len(m.getActiveTrack().Steps()) - 1) % stepsPerPage
		if m.activeTrackPage > 0 && remainingStepsInPage == 0 {
			m.activeTrackPage--
		}
		m.activeStep = 0
		m.seq.RemoveStep(m.activeTrack)
	}
}

func (m *mainModel) getActiveTrack() sequencer.Track {
	return m.seq.Tracks()[m.activeTrack]
}

func (m *mainModel) getActiveStep() sequencer.Step {
	return m.seq.Tracks()[m.activeTrack].Steps()[m.activeStep]
}
