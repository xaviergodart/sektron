package ui

import (
	"sektron/sequencer"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	refreshFrequency = 16 * time.Millisecond
)

type TickMsg time.Time

type mainModel struct {
	seq             sequencer.SequencerInterface
	width           int
	height          int
	pressedKey      *tea.KeyMsg
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

		// These keys should exit the program.
		case "ctrl+c", "esc":
			m.seq.Reset()
			return m, tea.Quit
		}

		switch {
		case key.Matches(msg, DefaultKeyMap.Step1):
			m.activeTrack = 0
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step2):
			m.activeTrack = 1
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step3):
			m.activeTrack = 2
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step4):
			m.activeTrack = 3
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step5):
			m.activeTrack = 4
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step6):
			m.activeTrack = 5
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step7):
			m.activeTrack = 6
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step8):
			m.activeTrack = 7
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step9):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step10):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step11):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step12):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step13):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step14):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step15):
			return m, nil
		case key.Matches(msg, DefaultKeyMap.Step16):
			return m, nil
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	return m.renderSequencer()
}
