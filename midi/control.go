package midi

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi/v2"
)

const (
	minCC      = 0
	maxCC      = 127
	minPitch   = -8192
	maxPitch   = 8192
	resetPitch = 0
)

type msgType uint8

const (
	controlChange msgType = iota
	programChange
	pitchBend
	afterTouch
)

type Controllable interface {
	Device() int
	Channel() uint8
	Controls() []*Control
	IsActiveControl(control uint8) bool
}

type Control struct {
	midi       Midi
	parent     Controllable
	msgType    msgType
	controller uint8
	value      int16
}

func NewControls(midi Midi, parent Controllable) []*Control {
	controls := []*Control{
		{
			midi:    midi,
			parent:  parent,
			msgType: programChange,
		},
		{
			midi:    midi,
			parent:  parent,
			msgType: pitchBend,
		},
		{
			midi:    midi,
			parent:  parent,
			msgType: afterTouch,
		},
	}

	for i := 0; i <= 127; i++ {
		controls = append(controls, &Control{
			midi:       midi,
			parent:     parent,
			msgType:    controlChange,
			controller: uint8(i),
		})
	}
	return controls
}

func (c Control) Value() int16 {
	return c.value
}

func (c Control) String() string {
	return fmt.Sprintf("%d", c.value)
}

func (c Control) Name() string {
	return gomidi.ControlChangeName[c.controller]
}

func (c *Control) Set(value int16) {
	switch c.msgType {
	case pitchBend:
		if value < minPitch || value > maxPitch {
			return
		}
	default:
		if value < minCC || value > maxCC {
			return
		}
	}
	shouldTrig := value != c.value
	c.value = value
	if shouldTrig {
		c.Trigger()
	}
}

func (c Control) Trigger() {
	switch c.msgType {
	case controlChange:
		c.midi.ControlChange(c.parent.Device(), c.parent.Channel(), c.controller, uint8(c.value))
	case programChange:
		c.midi.ProgramChange(c.parent.Device(), c.parent.Channel(), uint8(c.value))
	case pitchBend:
		c.midi.Pitchbend(c.parent.Device(), c.parent.Channel(), c.value)
	case afterTouch:
		c.midi.AfterTouch(c.parent.Device(), c.parent.Channel(), uint8(c.value))
	}
}
