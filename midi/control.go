package midi

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi/v2"
)

const (
	minCC = 0
	maxCC = 127
)

type msgType uint8

const (
	controlChange msgType = iota
	programChange
	pitchBend
	afterTouch
)

type Control struct {
	msgType    msgType
	controller uint8
	value      uint8
}

func NewControls() []*Control {
	controls := []*Control{
		{
			msgType: programChange,
		},
		{
			msgType: pitchBend,
		},
		{
			msgType: afterTouch,
		},
	}

	for i := 0; i <= 127; i++ {
		controls = append(controls, &Control{
			msgType:    controlChange,
			controller: uint8(i),
		})
	}
	return controls
}

func (c Control) Value() uint8 {
	return c.value
}

func (c Control) String() string {
	return fmt.Sprintf("%d", c.value)
}

func (c Control) Name() string {
	return gomidi.ControlChangeName[c.controller]
}

func (c *Control) Set(value uint8) {
	if value < minCC || value > maxCC {
		return
	}
	if value != c.value {
		c.Trigger()
	}
	c.value = value
}

func (c Control) Trigger() {
	return
}
