package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

var (
	stepKeys  = []string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"}
	trackKeys = []string{"A", "Z", "E", "R", "T", "Y", "U", "I", "Q", "S", "D", "F", "G", "H", "J", "K"}
	paramKeys = []string{"&", "é", "\"", "'", "(", "-", "è", "_", "ç", "à"}
)

type keyMap struct {
	TogglePlay key.Binding
	Mode       key.Binding

	Add    key.Binding
	Remove key.Binding

	StepsIndex map[string]int
	Steps      key.Binding

	TracksIndex   map[string]int
	Tracks        key.Binding
	TrackPageUp   key.Binding
	TrackPageDown key.Binding

	TempoUp       key.Binding
	TempoDown     key.Binding
	TempoFineUp   key.Binding
	TempoFineDown key.Binding

	ParamsIndex map[string]int
	Params      key.Binding

	ParamValueUp   key.Binding
	ParamValueDown key.Binding

	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TogglePlay, k.Mode, k.Add, k.Remove, k.TempoUp, k.TempoDown, k.TempoFineUp, k.TempoFineDown},
		{k.Steps, k.Tracks, k.TrackPageUp, k.TrackPageDown},
		{k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default key mapping.
func DefaultKeyMap() keyMap {
	km := keyMap{
		TogglePlay: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle play"),
		),
		Mode: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle mode (track, record)"),
		),
		Add: key.NewBinding(
			key.WithKeys("="),
			key.WithHelp("=", "add track|step"),
		),
		Remove: key.NewBinding(
			key.WithKeys(")"),
			key.WithHelp(")", "remove track|step"),
		),
		StepsIndex: map[string]int{},
		Steps: key.NewBinding(
			key.WithKeys(stepKeys...),
			key.WithHelp("a...i | q...k", "select track|step 1 to 16"),
		),
		TracksIndex: map[string]int{},
		Tracks: key.NewBinding(
			key.WithKeys(trackKeys...),
			key.WithHelp("A...I | Q...K", "toggle track 1 to 16"),
		),
		TrackPageUp: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "track page up"),
		),
		TrackPageDown: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "track page down"),
		),
		TempoUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("page up", "tempo up (1 bpm)"),
		),
		TempoDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("page down", "tempo down (1 bpm)"),
		),
		TempoFineUp: key.NewBinding(
			key.WithKeys("alt+pgup"),
			key.WithHelp("alt+page up", "tempo up (0.1 bpm)"),
		),
		TempoFineDown: key.NewBinding(
			key.WithKeys("alt+pgdown"),
			key.WithHelp("alt+page down", "tempo down (0.1 bpm)"),
		),
		ParamsIndex: map[string]int{},
		Params: key.NewBinding(
			key.WithKeys(paramKeys...),
			key.WithHelp(strings.Join(paramKeys, "/"), "select parameter 1 to 10"),
		),
		ParamValueUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "increase selected parameter value"),
		),
		ParamValueDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "decrease selected parameter value"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "quit"),
		),
	}
	for i, k := range stepKeys {
		km.StepsIndex[k] = i
	}
	for i, k := range trackKeys {
		km.TracksIndex[k] = i
	}
	for i, k := range paramKeys {
		km.ParamsIndex[k] = i
	}
	return km
}
