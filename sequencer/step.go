package sequencer

import (
	"math/rand"
	"sektron/midi"
	"strconv"
)

// Step contains a step state.
type Step interface {
	Track() *track
	Position() int
	IsActive() bool
	IsCurrentStep() bool
	Offset() int
	OffsetString() string
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

	// an offset relative to the first pulse on the step can be defined. It
	// delays the step trigger by x pulses, from 0 to 5.
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

// ChordString returns the string representation of the step chord.
func (s step) ChordString() string {
	return chordString(s.Chord())
}

// VelocityString returns the string representation of the step velocity.
func (s step) VelocityString() string {
	return velocityString(s.Velocity())
}

// LengthString returns the string representation of the step length.
func (s step) LengthString() string {
	return lengthString(s.Length())
}

// ProbabilityString returns the string representation of the step probability.
func (s step) ProbabilityString() string {
	return probabilityString(s.Probability())
}

// OffsetString returns the string representation of the step offset.
func (s step) OffsetString() string {
	return strconv.Itoa(s.offset)
}

// SetChord sets a new chord value.
func (s *step) SetChord(chord []uint8) {
	for _, note := range chord {
		if note < minChordNote || note > maxChordNote {
			return
		}
	}
	s.reset()
	s.chord = &chord
}

// SetLength sets a new length value.
func (s *step) SetLength(length int) {
	if length < minLength || length > maxLength {
		return
	}
	s.length = &length
}

// SetVelocity sets a new velocity value.
func (s *step) SetVelocity(velocity uint8) {
	if velocity < minVelocity || velocity > maxVelocity {
		return
	}
	s.velocity = &velocity
}

// SetProbability sets a new probability value.
func (s *step) SetProbability(probability int) {
	if probability < minProbability || probability > maxProbability {
		return
	}
	s.probability = &probability
}

// SetOffset sets a new offset value
func (s *step) SetOffset(offset int) {
	if offset < minOffset || offset > maxOffset {
		return
	}
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
	s.track.lastTriggeredStep = s.position
}

func (s step) skip() bool {
	return s.Probability() < 100 && rand.Intn(100) > s.Probability()
}

func (s step) startingPulse() int {
	return s.position*pulsesPerStep + s.offset
}

func (s step) endingPulse() int {
	return (s.startingPulse() + s.Length() - 1) % (pulsesPerStep * len(s.track.steps))
}

func (s step) isStartingPulse() bool {
	return s.track.pulse == s.startingPulse()
}

func (s step) isEndingPulse() bool {
	return s.track.pulse == s.endingPulse()
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
