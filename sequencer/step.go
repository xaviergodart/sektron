package sequencer

import (
	"math/rand"
	"sektron/midi"
)

type Step struct {
	midi   midi.MidiInterface
	track  *Track
	number int

	length      *int
	chord       *[]uint8
	velocity    *uint8
	probability *int
	offset      int
	active      bool
	triggered   bool
}

func (s Step) Track() *Track {
	return s.track
}

func (s Step) Number() int {
	return s.number
}

func (s Step) IsCurrentStep() bool {
	return s.number == s.track.CurrentStep()
}

func (s Step) Chord() []uint8 {
	if s.chord == nil {
		return s.track.chord
	}
	return *s.chord
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

func (s Step) IsActive() bool {
	return s.active
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
	return s.relativePulse() == s.Length()-1+s.offset
}

func (s *Step) trigger() {
	if !s.active || s.triggered || s.skip() {
		return
	}
	for _, note := range s.Chord() {
		s.midi.NoteOn(s.track.device, s.track.channel, note, s.Velocity())
	}
	s.triggered = true
}

func (s *Step) reset() {
	if !s.triggered {
		return
	}
	for _, note := range s.Chord() {
		s.midi.NoteOff(s.track.device, s.track.channel, note)
	}
	s.triggered = false
}
