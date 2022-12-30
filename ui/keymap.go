package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Step1  key.Binding
	Step2  key.Binding
	Step3  key.Binding
	Step4  key.Binding
	Step5  key.Binding
	Step6  key.Binding
	Step7  key.Binding
	Step8  key.Binding
	Step9  key.Binding
	Step10 key.Binding
	Step11 key.Binding
	Step12 key.Binding
	Step13 key.Binding
	Step14 key.Binding
	Step15 key.Binding
	Step16 key.Binding
}

var DefaultKeyMap = KeyMap{
	Step1: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "select track 1 | Rec mode: toggle step 1"),
	),
	Step2: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("z", "select track 2 | Rec mode: toggle step 2"),
	),
	Step3: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "select track 3 | Rec mode: toggle step 3"),
	),
	Step4: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "select track 4 | Rec mode: toggle step 4"),
	),
	Step5: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "select track 5 | Rec mode: toggle step 5"),
	),
	Step6: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "select track 6 | Rec mode: toggle step 6"),
	),
	Step7: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "select track 7 | Rec mode: toggle step 7"),
	),
	Step8: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "select track 8 | Rec mode: toggle step 8"),
	),
	Step9: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Rec mode: toggle step 9"),
	),
	Step10: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Rec mode: toggle step 10"),
	),
	Step11: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Rec mode: toggle step 11"),
	),
	Step12: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Rec mode: toggle step 12"),
	),
	Step13: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "Rec mode: toggle step 13"),
	),
	Step14: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Rec mode: toggle step 14"),
	),
	Step15: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "Rec mode: toggle step 15"),
	),
	Step16: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "Rec mode: toggle step 16"),
	),
}
