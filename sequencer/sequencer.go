// Package sequencer provides structures that holds the sequencer state, and
// ways to control that state and make it evolve over time.
//
// A sequencer instance is composed of:
//   - a clock that handle time events
//   - 1 to 16 tracks
//   - up to 64 steps per track
//
// Uppon creation, a new sequencer should receive an instrument that is allowed
// to play notes.
package sequencer

import (
	"math/rand"

	"sektron/instrument"
	"time"
)

const (
	defaultTempo         float64 = 120.0
	defaultTracks        int     = 8
	minTracks            int     = 1
	maxTracks            int     = 16
	defaultNote          uint8   = 60
	defaultVelocity      uint8   = 100
	defaultProbability   int     = 100
	defaultDevice        int     = 0
	defaultStepsPerTrack int     = 16
)

// Sequencer contains the sequencer state.
type Sequencer interface {
	TogglePlay()
	IsPlaying() bool
	AddTrack()
	RemoveTrack()
	Tracks() []*track
	ToggleTrack(track int)
	ToggleStep(track int, step int)
	Tempo() float64
	SetTempo(tempo float64)
	Reset()
}

type sequencer struct {
	instrument instrument.Instrument
	tracks     []*track
	clock      *clock

	// Holds the devices to which we should send the clock.
	// Useful for midi instrument.
	clockSend []int

	isPlaying bool
}

// New creates a new sequencer. It also creates new tracks and calls the
// start() method that starts the clock.
func New(instrument instrument.Instrument) *sequencer {
	// The randomizer will be used for step trigger probability.
	// Check step.go.
	rand.Seed(time.Now().UnixNano())

	seq := &sequencer{
		instrument: instrument,
		clockSend:  []int{defaultDevice},
		isPlaying:  false,
	}

	for i := 0; i < defaultTracks; i++ {
		seq.AddTrack()
	}

	// Let's start the clock right away.
	seq.start()
	return seq
}

// TogglePlay plays or stops the sequencer. When stopping, the sequencer resets
// the playhead to the first step and stops all the playing notes.
func (s *sequencer) TogglePlay() {
	s.isPlaying = !s.isPlaying
	if !s.isPlaying {
		s.Reset()
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
		pulse:       pulse,
		chord:       []uint8{defaultNote},
		length:      pulsesPerStep,
		velocity:    defaultVelocity,
		probability: defaultProbability,
		device:      defaultDevice,
		channel:     uint8(channel),
		active:      true,
	}

	var steps []*step
	for j := 0; j < defaultStepsPerTrack; j++ {
		steps = append(steps, &step{
			number:     j,
			instrument: s.instrument,
			track:      track,
			active:     false,
		})
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
}

func (s *sequencer) start() {
	// Each time the clock ticks, we call the sequencer tick method that
	// basically makes every track move forward in time.
	s.clock = newClock(defaultTempo, func() {
		s.tick()
	})
}

func (s *sequencer) tick() {
	// We send clock tick to the instrument in case it can sync with it.
	// Useful mainly for midi instrument.
	s.instrument.SendClock(s.clockSend)

	if !s.isPlaying {
		return
	}
	for _, track := range s.tracks {
		track.tick()
	}
}
