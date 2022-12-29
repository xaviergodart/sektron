package sequencer

import (
	"math/rand"

	"sektron/midi"
)

type Step struct {
	midi   *midi.Server
	track  *Track
	number int

	length      *int
	note        *uint8
	velocity    *uint8
	probability *int
	offset      int
	active      bool
	triggered   bool
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

func (s Step) Length() int {
	if s.length == nil {
		return s.track.length
	}
	return *s.length
}

func (s Step) Probability() int {
	if s.probability == nil {
		return s.track.probability
	}
	return *s.probability
}

func (s Step) skip() bool {
	return s.Probability() < 100 && rand.Intn(100) > s.Probability()
}

func (s Step) relativePulse() int {
	return s.track.pulse - (s.number * pulsesPerStep)
}

func (s Step) isStartingPulse() bool {
	return s.relativePulse() == s.offset
}

func (s Step) isEndingPulse() bool {
	return s.relativePulse() >= s.Length()+s.offset
}

func (s *Step) trigger() {
	if !s.active || s.triggered || s.skip() {
		return
	}
	s.midi.NoteOn(s.track.device, s.track.channel, s.Note(), s.Velocity())
	s.triggered = true
}

func (s *Step) reset() {
	if !s.triggered {
		return
	}
	s.midi.NoteOff(s.track.device, s.track.channel, s.Note())
	s.triggered = false
}
