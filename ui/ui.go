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

type RefreshTickMsg time.Time

type UI struct {
	seq             sequencer.SequencerInterface
	activeTrack     int
	activeTrackPage int
}

func New(seq sequencer.SequencerInterface) UI {
	return UI{
		seq:             seq,
		activeTrack:     0,
		activeTrackPage: 0,
	}
}

func refresh() tea.Cmd {
	return tea.Tick(refreshFrequency, func(t time.Time) tea.Msg {
		return RefreshTickMsg(t)
	})
}

func (u UI) Init() tea.Cmd {
	return refresh()
}

func (u UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case RefreshTickMsg:
		return u, refresh()

	case tea.KeyMsg:
		switch msg.String() {

		case " ":
			u.seq.TogglePlay()
			return u, nil

		// These keys should exit the program.
		case "ctrl+c", "esc":
			u.seq.Reset()
			return u, tea.Quit
		}

		switch {
		case key.Matches(msg, DefaultKeyMap.Step1):
			u.activeTrack = 0
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step2):
			u.activeTrack = 1
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step3):
			u.activeTrack = 2
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step4):
			u.activeTrack = 3
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step5):
			u.activeTrack = 4
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step6):
			u.activeTrack = 5
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step7):
			u.activeTrack = 6
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step8):
			u.activeTrack = 7
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step9):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step10):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step11):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step12):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step13):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step14):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step15):
			return u, nil
		case key.Matches(msg, DefaultKeyMap.Step16):
			return u, nil
		}
	}
	return u, nil
}

func (u UI) View() string {
	return u.ViewSequencer()
}
