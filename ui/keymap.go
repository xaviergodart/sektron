package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

var (
	stepKeys  = []string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"}
	paramKeys = []string{"&", "é", "\"", "'", "(", "-", "è", "_", "ç", "à"}
)

type KeyMap struct {
	StepsIndex map[string]int
	Steps      key.Binding

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
	for i, k := range paramKeys {
		km.ParamsIndex[k] = i
	}
	return km
}
