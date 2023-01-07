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

type KeyMap struct {
	StepsIndex map[string]int
	Steps      key.Binding

	TracksIndex   map[string]int
	Tracks        key.Binding
	TrackPageUp   key.Binding
	TrackPageDown key.Binding

	ParamsIndex map[string]int
	Params      key.Binding

	ParamValueUp   key.Binding
	ParamValueDown key.Binding
}

func DefaultKeyMap() KeyMap {
	km := KeyMap{
		StepsIndex: map[string]int{},
		Steps: key.NewBinding(
			key.WithKeys(stepKeys...),
			key.WithHelp(strings.Join(stepKeys, "/"), "track mode: select tracks 1 to 16 | record mode: toggle step 1 - 16"),
		),
		TracksIndex: map[string]int{},
		Tracks: key.NewBinding(
			key.WithKeys(trackKeys...),
			key.WithHelp(strings.Join(stepKeys, "/"), "toggle tracks 1 to 16"),
		),
		TrackPageUp: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "track page up"),
		),
		TrackPageDown: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "track page down"),
		),
		ParamsIndex: map[string]int{},
		Params: key.NewBinding(
			key.WithKeys(paramKeys...),
			key.WithHelp(strings.Join(paramKeys, "/"), "select parameters 1 to 10"),
		),
		ParamValueUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "increase selected parameter value"),
		),
		ParamValueDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "decrease selected parameter value"),
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
