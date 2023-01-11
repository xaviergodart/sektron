package sequencer

import (
	"math/rand"
	"sektron/instrument"
)

type StepInterface interface {
	Track() *track
	Number() int
	IsActive() bool
	IsCurrentStep() bool
	Chord() []uint8
	Velocity() uint8
	Length() int
	Probability() int
}

type step struct {
	instrument instrument.Instrument
	track      *track
	number     int

	length      *int
	chord       *[]uint8
	velocity    *uint8
	probability *int
	offset      int
	active      bool
	triggered   bool
}

func (s step) Track() *track {
	return s.track
}

func (s step) Number() int {
	return s.number
}

func (s step) IsActive() bool {
	return s.active
}

func (s step) IsCurrentStep() bool {
	return s.number == s.track.CurrentStep()
}

func (s step) Chord() []uint8 {
	if s.chord == nil {
		return s.track.chord
	}
	return *s.chord
}

func (s step) Velocity() uint8 {
	if s.velocity == nil {
		return s.track.velocity
	}
	return *s.velocity
}

func (s step) Length() int {
	if s.length == nil {
		return s.track.length
	}
	return *s.length
}

func (s step) Probability() int {
	if s.probability == nil {
		return s.track.probability
	}
	return *s.probability
}

func (s step) skip() bool {
	return s.Probability() < 100 && rand.Intn(100) > s.Probability()
}

func (s step) relativePulse() int {
	return s.track.pulse - (s.number * pulsesPerStep)
}

func (s step) isStartingPulse() bool {
	return s.relativePulse() == s.offset
}

func (s step) isEndingPulse() bool {
	return s.relativePulse() == s.Length()-1+s.offset
}

func (s *step) trigger() {
	if !s.active || s.triggered || s.skip() {
		return
	}
	for _, note := range s.Chord() {
		s.instrument.NoteOn(s.track.device, s.track.channel, note, s.Velocity())
	}
	s.triggered = true
}

func (s *step) reset() {
	if !s.triggered {
		return
	}
	for _, note := range s.Chord() {
		s.instrument.NoteOff(s.track.device, s.track.channel, note)
	}
	s.triggered = false
}
