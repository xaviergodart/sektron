package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

var (
	stepKeys        = []string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"}
	stepToggleKeys  = []string{"A", "Z", "E", "R", "T", "Y", "U", "I", "Q", "S", "D", "F", "G", "H", "J", "K"}
	trackKeys       = []string{"&", "é", "\"", "'", "(", "-", "è", "_", "ç", "à"}
	trackToggleKeys = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
)

type keyMap struct {
	Play key.Binding

	ParamMode   key.Binding
	PatternMode key.Binding

	AddTrack    key.Binding
	RemoveTrack key.Binding

	AddStep    key.Binding
	RemoveStep key.Binding

	StepIndex map[string]int
	Step      key.Binding

	StepToggleIndex map[string]int
	StepToggle      key.Binding

	TrackIndex map[string]int
	Track      key.Binding

	TrackToggleIndex map[string]int
	TrackToggle      key.Binding

	PageUp   key.Binding
	PageDown key.Binding

	TempoUp       key.Binding
	TempoDown     key.Binding
	TempoFineUp   key.Binding
	TempoFineDown key.Binding

	AddParam    key.Binding
	RemoveParam key.Binding

	Validate key.Binding
	Left     key.Binding
	Right    key.Binding
	Up       key.Binding
	Down     key.Binding

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
		{k.Play, k.ParamMode, k.PatternMode, k.AddTrack, k.RemoveTrack, k.AddStep, k.RemoveStep, k.TempoUp, k.TempoDown, k.TempoFineUp, k.TempoFineDown},
		{k.Step, k.StepToggle, k.Track, k.TrackToggle, k.PageUp, k.PageDown, k.AddParam, k.RemoveParam},
		{k.Validate, k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default key mapping.
func DefaultKeyMap() keyMap {
	km := keyMap{
		Play: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle play"),
		),
		ParamMode: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle parameter mode (track, record)"),
		),
		PatternMode: key.NewBinding(
			key.WithKeys("²"),
			key.WithHelp("²", "toggle pattern select mode"),
		),
		AddTrack: key.NewBinding(
			key.WithKeys("="),
			key.WithHelp("=", "add track"),
		),
		RemoveTrack: key.NewBinding(
			key.WithKeys(")"),
			key.WithHelp(")", "remove track"),
		),
		AddStep: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "add step"),
		),
		RemoveStep: key.NewBinding(
			key.WithKeys("°"),
			key.WithHelp("°", "remove track"),
		),
		StepIndex: map[string]int{},
		Step: key.NewBinding(
			key.WithKeys(stepKeys...),
			key.WithHelp(strings.Join(stepKeys, "/"), "select step|pattern 1 to 16"),
		),
		StepToggleIndex: map[string]int{},
		StepToggle: key.NewBinding(
			key.WithKeys(stepToggleKeys...),
			key.WithHelp(strings.Join(stepToggleKeys, "/"), "toggle step or chain pattern 1 to 16"),
		),
		TrackIndex: map[string]int{},
		Track: key.NewBinding(
			key.WithKeys(trackKeys...),
			key.WithHelp(strings.Join(trackKeys, "/"), "select track 1 to 10"),
		),
		TrackToggleIndex: map[string]int{},
		TrackToggle: key.NewBinding(
			key.WithKeys(trackToggleKeys...),
			key.WithHelp(strings.Join(trackToggleKeys, "/"), "toggle track 1 to 10"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "step|pattern page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "step|pattern page down"),
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
			key.WithKeys("ctrl+pgup"),
			key.WithHelp("ctrl+page up", "tempo up (0.1 bpm)"),
		),
		TempoFineDown: key.NewBinding(
			key.WithKeys("ctrl+pgdown"),
			key.WithHelp("ctrl+page down", "tempo down (0.1 bpm)"),
		),
		AddParam: key.NewBinding(
			key.WithKeys("ctrl+up"),
			key.WithHelp("ctrl+↑", "add midi control"),
		),
		RemoveParam: key.NewBinding(
			key.WithKeys("ctrl+down"),
			key.WithHelp("ctrl+↓", "remove midi control"),
		),
		Validate: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "validate selection"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "parameter select left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "parameter select left"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "increase selected parameter value"),
		),
		Down: key.NewBinding(
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
		km.StepIndex[k] = i
	}
	for i, k := range stepToggleKeys {
		km.StepToggleIndex[k] = i
	}
	for i, k := range trackKeys {
		km.TrackIndex[k] = i
	}
	for i, k := range trackToggleKeys {
		km.TrackToggleIndex[k] = i
	}
	return km
}
