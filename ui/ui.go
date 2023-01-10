package ui

import (
	"sektron/sequencer"
	"time"

	"github.com/charmbracelet/bubbles/help"
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
	refreshFrequency   = 16 * time.Millisecond // TODO: should be ok up to 50ms
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
	help            help.Model
}

func New(seq sequencer.SequencerInterface) mainModel {
	return mainModel{
		seq:             seq,
		keymap:          DefaultKeyMap(),
		activeTrack:     0,
		activeTrackPage: 0,
		activeParam:     0,
		help:            help.New(),
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
		m.help.Width = msg.Width
		return m, nil

	case TickMsg:
		return m, tick()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.TogglePlay):
			m.seq.TogglePlay()
			return m, nil

		case key.Matches(msg, m.keymap.Mode):
			if m.mode == trackMode {
				m.mode = recMode
			} else {
				m.mode = trackMode
			}
			return m, nil

		case key.Matches(msg, m.keymap.Add):
			m.seq.AddTrack()
			return m, nil

		case key.Matches(msg, m.keymap.Remove):
			if m.activeTrack > 0 && m.activeTrack == len(m.seq.Tracks())-1 {
				m.activeTrack--
			}
			m.seq.RemoveTrack()
			return m, nil

		case key.Matches(msg, m.keymap.Steps):
			m.stepPress(msg)
			return m, nil

		case key.Matches(msg, m.keymap.Tracks):
			number := m.keymap.TracksIndex[msg.String()]
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

		case key.Matches(msg, m.keymap.Params):
			m.activeParam = m.keymap.ParamsIndex[msg.String()]
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

func (m mainModel) trackPagesNb() int {
	pageNb := len(m.seq.Tracks()[m.activeTrack].Steps()) / stepsPerPage
	if len(m.seq.Tracks()[m.activeTrack].Steps())%stepsPerPage > 0 {
		pageNb++
	}
	return pageNb
}

func (m mainModel) View() string {
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderStatus(),
		m.renderSequencer(),
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
