package sequencer

import (
	"fmt"
	"sektron/midi"
)

// Track contains a track state.
type Track interface {
	Device() int
	DeviceString() string
	SetDevice(device int)
	Channel() uint8
	ChannelString() string
	SetChannel(channel uint8)
	Steps() []*step
	CurrentStep() int
	IsActive() bool
	IsCurrentStepActive() bool
	Parametrable
}

type track struct {
	midi  midi.Midi
	steps []*step

	// The pulse defines the current position of the playhead in the track.
	// Each time the clock ticks, we increment the pulse.
	// pulse ranges from 0 to len(steps) * pulsesPerStep (check clock.go).
	// Because each track can have a different number of steps, track pulses
	// are not always synchronized.
	pulse int

	// A track can be assigned to a specific midi device, channel and program.
	device  int
	channel uint8

	// A track can send multiple midi control changes. All possible midi
	// controls are initialized in the controls slice. But actual messages
	// will be sent only for those that are activated.
	controls       []midi.Control
	activeControls map[int]struct{}

	// Each track starts a goroutine to handle its pulse progression and step
	// triggering, by using the trig chan at each clock tick.
	// On track removal, we use the done chan to terminate the goroutine.
	trig chan struct{}
	done chan struct{}

	// An inactive track will progress like an active track, but will not
	// trigger any steps.
	active bool

	// We store the last triggered step of each track in order to reset it
	// if a new step is triggered. We avoid to steps being triggered at the same
	// time.
	lastTriggeredStep int

	// The next attributes defines the note parameters for the midi note on/off
	// messages and can be overriden per step (check step.go).
	//  - length defines for how long (pulse value) the note should be played
	//  - chord holds all the notes that should be played
	//  - velocity defines how loud a note should be played
	//  - probability defines the chances that the note will be played
	length      int
	chord       []uint8
	velocity    uint8
	probability int
}

// Steps returns all the track steps.
func (t track) Steps() []*step {
	return t.steps
}

// CurrentStep returns the step where the pulse is right now.
func (t track) CurrentStep() int {
	return t.pulse / pulsesPerStep
}

// IsActive returns true if the track is active.
func (t track) IsActive() bool {
	return t.active
}

// IsCurrentStepActive returns true if the current step is active.
func (t track) IsCurrentStepActive() bool {
	if !t.active || len(t.steps) <= t.CurrentStep() {
		return false
	}
	return t.steps[t.CurrentStep()].IsActive()
}

// Device returns the track device.
func (t track) Device() int {
	return t.device
}

// DeviceString returns the device name string.
func (t track) DeviceString() string {
	return t.midi.Devices()[t.device].String()
}

// Channel returns the track midi channel.
func (t track) Channel() uint8 {
	return t.channel
}

// ChannelString returns the midi channel string.
func (t track) ChannelString() string {
	return fmt.Sprintf("%d", t.channel+1)
}

// Control returns the midi control from the number.
func (t track) Control(nb int) midi.Control {
	return t.controls[nb]
}

// IsActiveControl checks if a given control is active.
func (t track) IsActiveControl(control int) bool {
	_, active := t.activeControls[control]
	return active
}

// Chord returns the track chord.
func (t track) Chord() []uint8 {
	return t.chord
}

// Velocity returns the track velocity.
func (t track) Velocity() uint8 {
	return t.velocity
}

// Length returns the track length.
func (t track) Length() int {
	return t.length
}

// Probability returns the track probability.
func (t track) Probability() int {
	return t.probability
}

// ChordString returns the string representation of the track chord.
func (t track) ChordString() string {
	return chordString(t.chord)
}

// VelocityString returns the string representation of the track velocity.
func (t track) VelocityString() string {
	return velocityString(t.velocity)
}

// LengthString returns the string representation of the track length.
func (t track) LengthString() string {
	return lengthString(t.length)
}

// ProbabilityString returns the string representation of the track probability.
func (t track) ProbabilityString() string {
	return probabilityString(t.probability)
}

// SetDevice selects a device.
func (t *track) SetDevice(device int) {
	if device < 0 || len(t.midi.Devices()) <= device {
		return
	}
	t.clear()
	t.device = device
}

// SetChannel sets the midi channel.
func (t *track) SetChannel(channel uint8) {
	if channel < minChannel || channel > maxChannel {
		return
	}
	t.clear()
	t.channel = channel
}

// SetControl sets the given midi control.
func (t *track) SetControl(nb int, value int16) {
	t.controls[nb].Set(value)
}

// SetChord sets a new chord value.
func (t *track) SetChord(chord []uint8) {
	for _, note := range chord {
		if note < minChordNote || note > maxChordNote {
			return
		}
	}
	t.clear()
	t.chord = chord
}

// SetLength sets a new length value.
func (t *track) SetLength(length int) {
	if length < minLength {
		return
	}
	// Infinite mode
	if length > maxLength {
		t.length = maxLength
		return
	}
	t.length = length
}

// SetVelocity sets a new velocity value.
func (t *track) SetVelocity(velocity uint8) {
	if velocity < minVelocity || velocity > maxVelocity {
		return
	}
	t.velocity = velocity
}

// SetProbability sets a new probability value.
func (t *track) SetProbability(probability int) {
	if probability < minProbability || probability > maxProbability {
		return
	}
	t.probability = probability
}

func (t *track) start() {
	t.trig = make(chan struct{})
	t.done = make(chan struct{})
	go func(track *track) {
		for {
			select {
			case <-track.trig:
				track.trigger()
			case <-track.done:
				return
			}
		}
	}(t)
}

// tick will trigger the track at each clock tick.
func (t *track) tick() {
	t.trig <- struct{}{}
}

func (t *track) close() {
	defer close(t.done)
	defer close(t.trig)
	t.done <- struct{}{}
}

func (t track) previousStep() *step {
	return t.steps[t.lastTriggeredStep]
}

// trigger goes over each steps and trigger them or stop them if we're at their
// starting or ending pulse. They are calculated relative to the pulse, using
// the length and offset parameters (check step.go)
func (t *track) trigger() {
	for _, step := range t.steps {
		if t.active && step.isStartingPulse() {
			// We reset the last triggered step to avoid 2 steps of the same
			// track being triggered at the same time.
			if step.active && !t.previousStep().isInfinite() {
				t.previousStep().reset()
			}

			step.trigger()
			continue
		}

		if step.isEndingPulse() && !step.isInfinite() {
			step.reset()
		}
	}

	t.pulse++

	// Go back to the beginning if we reach the end of the track.
	if t.pulse == pulsesPerStep*len(t.steps) {
		t.pulse = 0
	}
}

func (t track) isInfinite() bool {
	return t.length == maxLength
}

// sendControls sends track's active midi control messages.
func (t track) sendControls() {
	for c := range t.activeControls {
		t.Control(c).Send()
	}
}

// reset move back the pulse to the beginning, and stops all the already
// triggered steps.
func (t *track) reset() {
	t.pulse = 0
	t.lastTriggeredStep = 0
	t.clear()
}

func (t *track) clear() {
	for _, step := range t.steps {
		step.reset()
	}
}
