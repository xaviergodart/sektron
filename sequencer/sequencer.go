// Package sequencer provides structures that holds the sequencer state, and
// ways to control that state and make it evolve over time.
//
// A sequencer instance is composed of:
//   - a clock that handle time events
//   - 1 to 10 tracks
//   - up to 128 steps per track
//
// Uppon creation, a new sequencer should receive a midi instance that is
// allowed to play notes.
package sequencer

import (
	"math/rand"

	"sektron/filesystem"
	"sektron/midi"
	"time"
)

const (
	defaultTempo         float64 = 120.0
	defaultTracks        int     = 4
	minTracks            int     = 1
	maxTracks            int     = 10
	defaultNote          uint8   = 60
	defaultVelocity      uint8   = 100
	defaultProbability   int     = 100
	defaultDevice        int     = 0
	defaultStepsPerTrack int     = 16
	minSteps             int     = 1
	maxSteps             int     = 128
)

// Sequencer contains the sequencer state.
type Sequencer interface {
	TogglePlay()
	IsPlaying() bool
	Save()
	Load(pattern int)
	LoadNextInChain()
	Chain(pattern int)
	ChainNow(pattern int)
	FullChain() []int
	Patterns() []filesystem.Pattern
	ActivePattern() int
	AddTrack()
	RemoveTrack()
	Tracks() []*track
	ToggleTrack(track int)
	AddStep(track int)
	RemoveStep(track int)
	ToggleStep(track int, step int)
	Tempo() float64
	SetTempo(tempo float64)
	Reset()
}

type sequencer struct {
	midi  midi.Midi
	bank  filesystem.Bank
	chain []int

	tracks []*track
	clock  *clock

	// Holds the midi devices to which we should send the clock.
	clockSend []int

	isPlaying bool

	isFirstTick bool
}

// New creates a new sequencer. It also creates new tracks and calls the
// start() method that starts the clock.
func New(midi midi.Midi, bank filesystem.Bank) *sequencer {
	// The randomizer will be used for step trigger probability.
	// Check step.go.
	rand.Seed(time.Now().UnixNano())

	seq := &sequencer{
		midi:        midi,
		bank:        bank,
		clockSend:   []int{defaultDevice},
		isPlaying:   false,
		isFirstTick: false,
	}

	// Let's start the clock right away.
	seq.start()

	// Load the last active pattern from bank if available.
	// Or instanciate default number of tracks.
	seq.Load(seq.bank.Active)

	return seq
}

// TogglePlay plays or stops the sequencer. When stopping, the sequencer resets
// the playhead to the first step and stops all the playing notes.
func (s *sequencer) TogglePlay() {
	s.isPlaying = !s.isPlaying
	if !s.isPlaying {
		s.Reset()
	} else {
		s.isFirstTick = true
		s.sendControls()
	}
}

// IsPlaying returns the sequencer playing status.
func (s *sequencer) IsPlaying() bool {
	return s.isPlaying
}

// AddTrack() adds a new track to the sequencer with defaults values and steps.
// You can add up to 16 tracks. It also starts the track (check track.go).
func (s *sequencer) AddTrack() {
	if len(s.tracks) == maxTracks {
		return
	}
	pulse := 0
	if len(s.tracks) > 0 {
		pulse = s.tracks[0].pulse
	}
	channel := len(s.tracks)
	track := &track{
		midi:                  s.midi,
		pulse:                 pulse,
		chord:                 []uint8{defaultNote},
		length:                pulsesPerStep,
		velocity:              defaultVelocity,
		probability:           defaultProbability,
		device:                defaultDevice,
		channel:               uint8(channel),
		activeControls:        map[int]struct{}{},
		lastSentControlValues: map[int]int16{},
		active:                true,
	}
	track.controls = midi.NewControls(s.midi, track)

	var steps []*step
	for j := 0; j < defaultStepsPerTrack; j++ {
		step := &step{
			position: j,
			midi:     s.midi,
			track:    track,
			active:   false,
			controls: map[int]*midi.Control{},
		}
		steps = append(steps, step)
	}

	track.steps = steps
	track.start()
	s.tracks = append(s.tracks, track)
}

// RemoveTrack removes the last track of the sequencer tracks. The first track
// can't be removed.
func (s *sequencer) RemoveTrack() {
	if len(s.tracks) == minTracks {
		return
	}
	s.tracks[len(s.tracks)-1].close()
	s.tracks = s.tracks[:len(s.tracks)-1]
}

// Tracks returns all the sequencer tracks.
func (s *sequencer) Tracks() []*track {
	return s.tracks
}

// AddStep adds a new step to the given track with default values.
// You can add up to 128 steps.
func (s *sequencer) AddStep(track int) {
	t := s.tracks[track]
	if len(t.steps) == maxSteps {
		return
	}
	t.steps = append(
		t.steps,
		&step{
			position: len(t.steps),
			midi:     s.midi,
			track:    t,
			active:   false,
		},
	)
}

// RemoveStep removes the last step of the given track. The first step
// can't be removed.
func (s *sequencer) RemoveStep(track int) {
	t := s.tracks[track]
	if len(t.steps) == minSteps {
		return
	}
	if t.lastTriggeredStep == len(t.steps)-1 {
		t.lastTriggeredStep = 0
	}
	t.steps[len(t.steps)-1].reset()
	t.steps = t.steps[:len(t.steps)-1]
	if t.pulse >= len(t.steps)*pulsesPerStep-1 {
		t.pulse = 0
	}
}

// Tempo returns the sequencer tempo.
func (s *sequencer) Tempo() float64 {
	return s.clock.tempo
}

// SetTempo allows to set the clock to a new tempo.
func (s *sequencer) SetTempo(tempo float64) {
	s.clock.setTempo(tempo)
}

// Reset resets all sequencer tracks (check track.go)
func (s *sequencer) Reset() {
	for _, track := range s.tracks {
		track.reset()
	}
}

// ToggleTrack activates or desactivates a specific track.
func (s *sequencer) ToggleTrack(track int) {
	if len(s.tracks) <= track {
		return
	}
	s.tracks[track].active = !s.tracks[track].active
}

// ToggleStep activates or desactivates a specific step of a given track.
func (s *sequencer) ToggleStep(track int, step int) {
	if len(s.tracks[track].steps) <= step {
		return
	}
	s.tracks[track].steps[step].active = !s.tracks[track].steps[step].active
	s.tracks[track].steps[step].clearParameters()
}

func (s *sequencer) start() {
	// Each time the clock ticks, we call the sequencer tick method that
	// basically makes every track move forward in time.
	s.clock = newClock(defaultTempo, func() {
		s.tick()
	})
}

func (s *sequencer) tick() {
	// We send clock tick to the midi devices.
	// TODO: make clock send configurable
	// s.midi.SendClock(s.clockSend)

	if !s.isPlaying {
		return
	}

	// Load first pattern in chain if chain not empty.
	if !s.isFirstTick && s.tracks[0].pulse == 0 {
		s.LoadNextInChain()
	}

	for _, track := range s.tracks {
		track.tick()
	}

	s.isFirstTick = false
}

// sendControls sends all track's active midi control messages.
func (s sequencer) sendControls() {
	for _, track := range s.tracks {
		track.sendControls()
	}
}
