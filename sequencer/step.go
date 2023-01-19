package sequencer

import (
	"math/rand"
	"sektron/midi"
)

// Step contains a step state.
type Step interface {
	Track() *track
	Position() int
	IsActive() bool
	IsCurrentStep() bool
	Offset() int
	SetOffset(offset int)
	Parametrable
}

type step struct {
	midi     midi.Midi
	track    *track
	position int

	// An inactive step will progress like an active step, but will not
	// trigger any notes.
	active bool

	// Once a step has been triggered, we prevent it from happening again.
	triggered bool

	// The next attributes defines the note parameters for the midi device.
	// If nil, we should use the default ones defined at track level (see
	// track.go)
	//  - length defines for how long (pulse value) the note should be played
	//  - chord holds all the notes that should be played
	//  - velocity defines how loud a note should be played
	//  - probability defines the chances that the note will be played
	length      *int
	chord       *[]uint8
	velocity    *uint8
	probability *int

	// an offset relative to the first pulse on the step can be defined, either
	// positive or negative. It allows microtimed triggers of the note.
	// 0 by default.
	offset int
}

// Track returns the parent track of the step.
func (s step) Track() *track {
	return s.track
}

// Position returns the position of the step in the track.
func (s step) Position() int {
	return s.position
}

// IsActive returns true if the step is active.
func (s step) IsActive() bool {
	return s.active
}

// IsCurrentStep returns true if the track pulse is on the current step.
func (s step) IsCurrentStep() bool {
	return s.position == s.track.CurrentStep()
}

// Chord returns the current step chord, or the one defined on the track if
// nil.
func (s step) Chord() []uint8 {
	if s.chord == nil {
		return s.track.chord
	}
	return *s.chord
}

// Velocity returns the current step velocity, or the one defined on the track
// if nil.
func (s step) Velocity() uint8 {
	if s.velocity == nil {
		return s.track.velocity
	}
	return *s.velocity
}

// Length returns the current step length, or the one defined on the track if
// nil.
func (s step) Length() int {
	if s.length == nil {
		return s.track.length
	}
	return *s.length
}

// Probability returns the current step probability, or the one defined on the
// track if nil.
func (s step) Probability() int {
	if s.probability == nil {
		return s.track.probability
	}
	return *s.probability
}

// Offset returns the current step offset value.
func (s step) Offset() int {
	return s.offset
}

// SetChord sets a new chord value.
func (s *step) SetChord(chord []uint8) {
	s.reset()
	s.chord = &chord
}

// SetLength sets a new length value.
func (s *step) SetLength(length int) {
	s.length = &length
}

// SetVelocity sets a new velocity value.
func (s *step) SetVelocity(velocity uint8) {
	s.velocity = &velocity
}

// SetProbability sets a new probability value.
func (s *step) SetProbability(probability int) {
	s.probability = &probability
}

// SetOffset sets a new offset value
func (s *step) SetOffset(offset int) {
	s.offset = offset
}

// Here we send the note on signal to the device if all the conditions are
// met. And we flag the step as triggered.
func (s *step) trigger() {
	if !s.active || s.triggered || s.skip() {
		return
	}
	for _, note := range s.Chord() {
		s.midi.NoteOn(s.track.device, s.track.channel, note, s.Velocity())
	}
	s.triggered = true
}

func (s step) skip() bool {
	return s.Probability() < 100 && rand.Intn(100) > s.Probability()
}

// relativePulse returns the step pulse relative to the track one. If negative,
// the track pulse is before the step. If zero, the track pulse is at the
// first step pulse. It allows starting and ending pulse calculation.
func (s step) relativePulse() int {
	return s.track.pulse - (s.position * pulsesPerStep)
}

// Starting and ending pulse are 0 and step length-1 relative to the track
// pulse. We can offset these to allow microtimed events.
func (s step) isStartingPulse() bool {
	return s.relativePulse() == s.offset
}

func (s step) isEndingPulse() bool {
	return s.relativePulse() == s.Length()-1+s.offset
}

// reset stops all the notes in the chord if the step has been triggered.
func (s *step) reset() {
	if !s.triggered {
		return
	}
	for _, note := range s.Chord() {
		s.midi.NoteOff(s.track.device, s.track.channel, note)
	}
	s.triggered = false
}
