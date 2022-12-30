package sequencer

import (
	"math/rand"

	"sektron/midi"
	"time"
)

const (
	defaultTempo         float64 = 120.0
	defaultNote          uint8   = 60
	defaultVelocity      uint8   = 100
	defaultProbability   int     = 100
	defaultDevice        int     = 0
	defaultStepsPerTrack int     = 16
)

type Sequencer struct {
	midi      *midi.Server
	tracks    []*Track
	clock     *Clock
	isPlaying bool
}

func New(midi *midi.Server) *Sequencer {
	rand.Seed(time.Now().UnixNano())
	var tracks []*Track
	for i := 0; i <= 1; i++ {
		var steps []*Step
		note := defaultNote + uint8(i*12) + uint8(i*5)
		track := &Track{
			pulse:       0,
			chord:       []uint8{note, note + 5},
			length:      pulsesPerStep,
			velocity:    defaultVelocity,
			probability: defaultProbability,
			device:      defaultDevice,
			channel:     uint8(i),
			active:      true,
		}
		for j := 0; j < defaultStepsPerTrack; j++ {
			steps = append(steps, &Step{
				number: j,
				midi:   midi,
				track:  track,
				active: j%4 == 0,
				offset: i * pulsesPerStep,
			})
		}
		track.steps = steps
		tracks = append(tracks, track)
	}
	return &Sequencer{
		midi:      midi,
		tracks:    tracks,
		isPlaying: false,
	}
}

func (s *Sequencer) Start() {
	for _, track := range s.tracks {
		track.Start()
	}
	s.clock = NewClock(defaultTempo, func() {
		s.Pulse()
	})
}

func (s *Sequencer) Tracks() []*Track {
	return s.tracks
}

func (s *Sequencer) TogglePlay() {
	s.isPlaying = !s.isPlaying
	if !s.isPlaying {
		s.Reset()
	}
}

func (s *Sequencer) Pulse() {
	s.midi.SendClock()
	if !s.isPlaying {
		return
	}
	for _, track := range s.tracks {
		track.Pulse()
	}
}

func (s *Sequencer) Reset() {
	for _, track := range s.tracks {
		track.reset()
	}
}
