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

type SequencerInterface interface {
	Reset()
	TogglePlay()
	Tracks() []*Track
}

type Sequencer struct {
	midi          midi.MidiInterface
	midiClockSend []int
	tracks        []*Track
	clock         *Clock
	isPlaying     bool
}

func New(midi midi.MidiInterface) *Sequencer {
	rand.Seed(time.Now().UnixNano())
	var tracks []*Track
	for i := 0; i <= 7; i++ {
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
	seq := &Sequencer{
		midi:          midi,
		midiClockSend: []int{defaultDevice},
		tracks:        tracks,
		isPlaying:     false,
	}
	seq.start()
	return seq
}

func (s *Sequencer) start() {
	for _, track := range s.tracks {
		track.start()
	}
	s.clock = NewClock(defaultTempo, func() {
		s.tick()
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

func (s *Sequencer) Reset() {
	for _, track := range s.tracks {
		track.reset()
	}
}

func (s *Sequencer) tick() {
	s.midi.SendClock(s.midiClockSend)
	if !s.isPlaying {
		return
	}
	for _, track := range s.tracks {
		track.tick()
	}
}
