package ui

import (
	"strings"

	"sektron/filesystem"

	"github.com/charmbracelet/bubbles/key"
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

	TempoUp   key.Binding
	TempoDown key.Binding

	AddParam    key.Binding
	RemoveParam key.Binding
	CopyParams  key.Binding
	PasteParams key.Binding

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
		{k.Play, k.ParamMode, k.PatternMode, k.AddTrack, k.RemoveTrack, k.AddStep, k.RemoveStep, k.TempoUp, k.TempoDown},
		{k.Step, k.StepToggle, k.Track, k.TrackToggle, k.PageUp, k.PageDown, k.AddParam, k.RemoveParam},
		{k.Validate, k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit},
	}
}

// newKeyMap returns the default key mapping.
func newKeyMap(keys filesystem.KeyMap) keyMap {
	km := keyMap{
		Play: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle play"),
		),
		ParamMode: key.NewBinding(
			key.WithKeys(keys.ParamMode),
			key.WithHelp(keys.ParamMode, "toggle parameter mode (track, step)"),
		),
		PatternMode: key.NewBinding(
			key.WithKeys(keys.PatternMode),
			key.WithHelp(keys.PatternMode, "toggle pattern select mode"),
		),
		AddTrack: key.NewBinding(
			key.WithKeys(keys.AddTrack),
			key.WithHelp(keys.AddTrack, "add track"),
		),
		RemoveTrack: key.NewBinding(
			key.WithKeys(keys.RemoveTrack),
			key.WithHelp(keys.RemoveTrack, "remove track"),
		),
		AddStep: key.NewBinding(
			key.WithKeys(keys.AddStep),
			key.WithHelp(keys.AddStep, "add step"),
		),
		RemoveStep: key.NewBinding(
			key.WithKeys(keys.RemoveStep),
			key.WithHelp(keys.RemoveStep, "remove step"),
		),
		StepIndex: map[string]int{},
		Step: key.NewBinding(
			key.WithKeys(keys.Steps[:]...),
			key.WithHelp(strings.Join(keys.Steps[:], "/"), "select step|pattern 1 to 16"),
		),
		StepToggleIndex: map[string]int{},
		StepToggle: key.NewBinding(
			key.WithKeys(keys.StepsToggle[:]...),
			key.WithHelp(strings.Join(keys.StepsToggle[:], "/"), "toggle step or chain pattern 1 to 16"),
		),
		TrackIndex: map[string]int{},
		Track: key.NewBinding(
			key.WithKeys(keys.Tracks[:]...),
			key.WithHelp(strings.Join(keys.Tracks[:], "/"), "select track 1 to 10"),
		),
		TrackToggleIndex: map[string]int{},
		TrackToggle: key.NewBinding(
			key.WithKeys(keys.TracksToggle[:]...),
			key.WithHelp(strings.Join(keys.TracksToggle[:], "/"), "toggle track 1 to 10"),
		),
		PageUp: key.NewBinding(
			key.WithKeys(keys.PageUp),
			key.WithHelp(keys.PageUp, "step|pattern page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys(keys.PageDown),
			key.WithHelp(keys.PageDown, "step|pattern page down"),
		),
		TempoUp: key.NewBinding(
			key.WithKeys(keys.TempoUp),
			key.WithHelp(keys.TempoUp, "tempo up (1 bpm)"),
		),
		TempoDown: key.NewBinding(
			key.WithKeys(keys.TempoDown),
			key.WithHelp(keys.TempoDown, "tempo down (1 bpm)"),
		),
		AddParam: key.NewBinding(
			key.WithKeys(keys.AddParam),
			key.WithHelp(keys.AddParam, "add midi control"),
		),
		RemoveParam: key.NewBinding(
			key.WithKeys(keys.RemoveParam),
			key.WithHelp(keys.RemoveParam, "remove midi control"),
		),
		CopyParams: key.NewBinding(
			key.WithKeys(keys.CopyParams),
			key.WithHelp(keys.CopyParams, "copy active step parameters"),
		),
		PasteParams: key.NewBinding(
			key.WithKeys(keys.PasteParams),
			key.WithHelp(keys.PasteParams, "paste copied parameters into active step"),
		),
		Validate: key.NewBinding(
			key.WithKeys(keys.Validate),
			key.WithHelp(keys.Validate, "validate selection"),
		),
		Left: key.NewBinding(
			key.WithKeys(keys.Left),
			key.WithHelp(keys.Left, "select parameter left"),
		),
		Right: key.NewBinding(
			key.WithKeys(keys.Right),
			key.WithHelp(keys.Right, "select parameter right"),
		),
		Up: key.NewBinding(
			key.WithKeys(keys.Up),
			key.WithHelp(keys.Up, "increase selected parameter value"),
		),
		Down: key.NewBinding(
			key.WithKeys(keys.Down),
			key.WithHelp(keys.Down, "decrease selected parameter value"),
		),
		Help: key.NewBinding(
			key.WithKeys(keys.Help),
			key.WithHelp(keys.Help, "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+q", "esc"),
			key.WithHelp("ctrl+q/esc", "quit"),
		),
	}
	for i, k := range keys.Steps {
		km.StepIndex[k] = i
	}
	for i, k := range keys.StepsToggle {
		km.StepToggleIndex[k] = i
	}
	for i, k := range keys.Tracks {
		km.TrackIndex[k] = i
	}
	for i, k := range keys.TracksToggle {
		km.TrackToggleIndex[k] = i
	}
	return km
}
