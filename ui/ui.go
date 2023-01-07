package ui

import (
	"sektron/sequencer"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time
type mode uint8

const (
	trackMode mode = iota
	recMode
)

const (
	sektron            = "SEKTRON"
	refreshFrequency   = 16 * time.Millisecond
	mainViewSideMargin = 2
)

type mainModel struct {
	seq             sequencer.SequencerInterface
	keymap          KeyMap
	width           int
	height          int
	mode            mode
	activeTrack     int
	activeTrackPage int
	activeParam     int
}

func New(seq sequencer.SequencerInterface) mainModel {
	return mainModel{
		seq:             seq,
		keymap:          DefaultKeyMap(),
		activeTrack:     0,
		activeTrackPage: 0,
		activeParam:     0,
	}
}

func tick() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return TickMsg(t)
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
		return m, nil

	case TickMsg:
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {

		case " ":
			m.seq.TogglePlay()
			return m, nil

		case "tab":
			if m.mode == trackMode {
				m.mode = recMode
			} else {
				m.mode = trackMode
			}
			return m, nil

		// These keys should exit the program.
		case "ctrl+c", "esc":
			m.seq.Reset()
			return m, tea.Quit
		}

		switch {
		case key.Matches(msg, m.keymap.Steps):
			m.stepPress(msg)
			return m, nil

		case key.Matches(msg, m.keymap.Tracks):
			number := m.keymap.TracksIndex[msg.String()]
			m.seq.ToggleTrack(number)
			return m, nil

		case key.Matches(msg, m.keymap.TrackPageUp):
			m.trackPagePress()
			return m, nil

		case key.Matches(msg, m.keymap.Params):
			m.activeParam = m.keymap.ParamsIndex[msg.String()]
			return m, nil
		}
	}
	return m, nil
}

func (m *mainModel) stepPress(msg tea.KeyMsg) {
	number := m.keymap.StepsIndex[msg.String()]
	switch m.mode {
	case trackMode:
		if number >= len(m.seq.Tracks()) {
			return
		}
		m.activeTrack = number
	case recMode:
		m.seq.ToggleStep(m.activeTrack, number+(m.activeTrackPage*stepsPerPage))
	}
}

func (m *mainModel) trackPagePress() {
	pageNb := len(m.seq.Tracks()[m.activeTrack].Steps()) / stepsPerPage
	if len(m.seq.Tracks()[m.activeTrack].Steps())%stepsPerPage > 0 {
		pageNb++
	}
	m.activeTrackPage = (m.activeTrackPage + 1) % pageNb
}

func (m mainModel) View() string {
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderStatus(),
		m.renderSequencer(),
	)

	// Cleanup gibber
	cleanup := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - lipgloss.Height(mainView)).
		Render("")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
		cleanup,
	)
}
