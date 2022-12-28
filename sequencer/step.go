package sequencer

import (
	"sektron/midi"
)

type step struct {
	midi      *midi.Server
	note      uint8
	triggered bool
}

func (s *step) trigger() {
	if s.triggered {
		return
	}
	// TODO: make device and channel configurable
	s.midi.NoteOn(0, 0, s.note, 100)
	s.triggered = true
}

func (s *step) reset() {
	if s.triggered {
		s.midi.NoteOff(0, 0, s.note)
	}
	s.triggered = false
}
