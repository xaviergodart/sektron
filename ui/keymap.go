package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Steps key.Binding
}

var stepKeys = []string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"}
var stepIndex = map[string]int{}

var DefaultKeyMap = KeyMap{
	Steps: key.NewBinding(
		key.WithKeys(stepKeys...),
		key.WithHelp(strings.Join(stepKeys, "/"), "track mode: select tracks 1 to 16 | record mode: toggle step 1 - 16"),
	),
}

func initKeyMap() {
	for i, k := range stepKeys {
		stepIndex[k] = i
	}
}
