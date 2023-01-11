package sequencer

import (
	"math/rand"

	"sektron/instrument"
	"time"
)

const (
	defaultTempo         float64 = 120.0
	defaultNote          uint8   = 60
	defaultVelocity      uint8   = 100
	defaultProbability   int     = 100
	defaultDevice        int     = 0
	defaultStepsPerTrack int     = 16

	minTracks     int = 1
	maxTracks     int = 16
	defaultTracks int = 8
)

type SequencerInterface interface {
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
	clockSend  []int
	tracks     []*track
	clock      *clock
	isPlaying  bool
}

func New(instrument instrument.Instrument) *sequencer {
	rand.Seed(time.Now().UnixNano())
	seq := &sequencer{
		instrument: instrument,
		clockSend:  []int{defaultDevice},
		isPlaying:  false,
	}

	for i := 0; i < defaultTracks; i++ {
		seq.AddTrack()
	}

	seq.start()
	return seq
}

func (s *sequencer) TogglePlay() {
	s.isPlaying = !s.isPlaying
	if !s.isPlaying {
		s.Reset()
	}
}

func (s *sequencer) IsPlaying() bool {
	return s.isPlaying
}

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

func (s *sequencer) RemoveTrack() {
	if len(s.tracks) == minTracks {
		return
	}
	s.tracks[len(s.tracks)-1].close()
	s.tracks = s.tracks[:len(s.tracks)-1]
}

func (s *sequencer) Tracks() []*track {
	return s.tracks
}

func (s *sequencer) Tempo() float64 {
	return s.clock.tempo
}

func (s *sequencer) SetTempo(tempo float64) {
	s.clock.setTempo(tempo)
}

func (s *sequencer) Reset() {
	for _, track := range s.tracks {
		track.reset()
	}
}

func (s *sequencer) ToggleTrack(track int) {
	if len(s.tracks) <= track {
		return
	}
	s.tracks[track].active = !s.tracks[track].active
}

func (s *sequencer) ToggleStep(track int, step int) {
	if len(s.tracks[track].steps) <= step {
		return
	}
	s.tracks[track].steps[step].active = !s.tracks[track].steps[step].active
}

func (s *sequencer) start() {
	s.clock = newClock(defaultTempo, func() {
		s.tick()
	})
}

func (s *sequencer) tick() {
	s.instrument.SendClock(s.clockSend)
	if !s.isPlaying {
		return
	}
	for _, track := range s.tracks {
		track.tick()
	}
}
