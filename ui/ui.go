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
	refreshFrequency = 16 * time.Millisecond
)

type mainModel struct {
	seq             sequencer.SequencerInterface
	width           int
	height          int
	pressedKey      *tea.KeyMsg
	mode            mode
	activeTrack     int
	activeTrackPage int
}

func New(seq sequencer.SequencerInterface) mainModel {
	return mainModel{
		seq:             seq,
		pressedKey:      nil,
		activeTrack:     0,
		activeTrackPage: 0,
	}
}

func tick() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m mainModel) Init() tea.Cmd {
	initKeyMap()
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
		m.pressedKey = &msg

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
		case key.Matches(msg, DefaultKeyMap.Steps):
			m.stepPress(msg)
			return m, nil
		}
	}
	return m, nil
}

func (m *mainModel) stepPress(msg tea.KeyMsg) {
	number := stepIndex[msg.String()]
	switch m.mode {
	case trackMode:
		if number >= 8 {
			return
		}
		m.activeTrack = number
	case recMode:
		m.seq.ToggleStep(m.activeTrack, number)
	}
}

func (m mainModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderSequencer(),
		m.renderStatus(),
	)
}
