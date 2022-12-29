package sequencer

import (
	"sektron/midi"
)

type Step struct {
	midi  *midi.Server
	track *Track

	note      *uint8
	velocity  *uint8
	active    bool
	triggered bool
}

func (s Step) Note() uint8 {
	if s.note == nil {
		return s.track.note
	}
	return *s.note
}

func (s Step) Velocity() uint8 {
	if s.velocity == nil {
		return s.track.velocity
	}
	return *s.velocity
}

func (s *Step) trigger() {
	if !s.active || s.triggered {
		return
	}
	s.midi.NoteOn(s.track.device, s.track.channel, s.Note(), s.Velocity())
	s.triggered = true
}

func (s *Step) reset() {
	if s.triggered {
		s.midi.NoteOff(s.track.device, s.track.channel, s.Note())
	}
	s.triggered = false
}
