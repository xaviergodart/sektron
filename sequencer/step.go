package sequencer

import (
	"fmt"
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

	// A step can override multiple midi control from the track.
	controls map[int]*midi.Control

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

// Device returns the step device.
func (s step) Device() int {
	return s.track.device
}

// Channel returns the step midi channel.
func (s step) Channel() uint8 {
	return s.track.channel
}

// Control returns the midi control from the number.
func (s step) Control(nb int) midi.Control {
	control, ok := s.controls[nb]
	if !ok {
		return s.track.controls[nb]
	}
	return *control
}

// IsActive returns true if the step is active.
func (s step) IsActive() bool {
	return s.active
}

// IsCurrentStep returns true if the track pulse is on the current step.
func (s step) IsCurrentStep() bool {
	return s.position == s.track.CurrentStep()
}

// IsActiveControl checks if a given control is active.
func (s step) IsActiveControl(control int) bool {
	_, active := s.track.activeControls[control]
	return active
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
	if s.offset == 0 {
		return fmt.Sprintf("%d", s.offset)
	}
	return fmt.Sprintf("+%d", s.offset)
}

// SetControl sets the given midi control.
func (s *step) SetControl(nb int, value int16) {
	_, ok := s.controls[nb]
	if !ok {
		control := s.track.controls[nb]
		s.controls[nb] = &control
	}
	s.controls[nb].Set(value)
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
	if length < minLength {
		return
	}
	// Infinite mode
	if length > maxLength {
		length = maxLength
		s.length = &length
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
	s.sendControls()
	for _, note := range s.Chord() {
		s.midi.NoteOn(s.track.device, s.track.channel, note, s.Velocity())
	}
	s.triggered = true
	s.track.lastTriggeredStep = s.position
}

// sendControls sends midi control messages if there step value are
// different from the previous step, to avoid sending the same messages
// multiple times.
func (s step) sendControls() {
	for c := range s.track.activeControls {
		// TODO: maybe this could be improved.
		// Small bug: if only the first step is activated,
		// controls are sent everytime.
		if s.isSameStepPlayed() && s.Control(c).Value() == s.track.Control(c).Value() {
			continue
		} else if !s.isSameStepPlayed() && s.Control(c).Value() == s.track.previousStep().Control(c).Value() {
			continue
		}
		s.Control(c).Send()
	}
}

func (s step) skip() bool {
	return s.Probability() < 100 && rand.Intn(100) > s.Probability()
}

func (s step) isSameStepPlayed() bool {
	return s.track.lastTriggeredStep == s.Position()
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

func (s step) isInfinite() bool {
	if s.length == nil {
		return s.track.isInfinite()
	}
	return *s.length == maxLength
}

func (s *step) clearParameters() {
	s.reset()
	s.length = nil
	s.chord = nil
	s.velocity = nil
	s.probability = nil
	s.offset = 0
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
