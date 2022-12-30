package ui

import (
	"sektron/sequencer"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	refreshFrequency = 16 * time.Millisecond
)

type RefreshTickMsg time.Time

type UI struct {
	seq *sequencer.Sequencer
}

func New(seq *sequencer.Sequencer) UI {
	return UI{
		seq: seq,
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
		case "ctrl+c", "q":
			u.seq.Reset()
			return u, tea.Quit
		}
	}
	return u, nil
}

func (u UI) View() string {
	return u.ViewSequencer()
}
