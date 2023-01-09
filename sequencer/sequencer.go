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

	minTracks     int = 1
	maxTracks     int = 16
	defaultTracks int = 8
)

type SequencerInterface interface {
	Reset()
	TogglePlay()
	IsPlaying() bool
	AddTrack()
	RemoveTrack()
	Tracks() []*Track
	ToggleTrack(track int)
	ToggleStep(track int, step int)
	Tempo() float64
	SetTempo(tempo float64)
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
	seq := &Sequencer{
		midi:          midi,
		midiClockSend: []int{defaultDevice},
		isPlaying:     false,
	}

	for i := 0; i < defaultTracks; i++ {
		seq.AddTrack()
	}

	seq.start()
	return seq
}

func (s *Sequencer) start() {
	s.clock = NewClock(defaultTempo, func() {
		s.tick()
	})
}

func (s *Sequencer) AddTrack() {
	if len(s.tracks) == maxTracks {
		return
	}
	channel := len(s.tracks)
	track := &Track{
		chord:       []uint8{defaultNote},
		length:      pulsesPerStep,
		velocity:    defaultVelocity,
		probability: defaultProbability,
		device:      defaultDevice,
		channel:     uint8(channel),
		active:      true,
	}

	var steps []*Step
	for j := 0; j < defaultStepsPerTrack; j++ {
		steps = append(steps, &Step{
			number: j,
			midi:   s.midi,
			track:  track,
			active: false,
		})
	}

	track.steps = steps
	track.start()
	s.tracks = append(s.tracks, track)
}

func (s *Sequencer) RemoveTrack() {
	if len(s.tracks) == minTracks {
		return
	}
	s.tracks = s.tracks[:len(s.tracks)-1]
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

func (s *Sequencer) Tempo() float64 {
	return s.clock.tempo
}

func (s *Sequencer) SetTempo(tempo float64) {
	s.clock.SetTempo(tempo)
}

func (s *Sequencer) IsPlaying() bool {
	return s.isPlaying
}

func (s *Sequencer) Reset() {
	for _, track := range s.tracks {
		track.clear()
	}
}

func (s *Sequencer) ToggleTrack(track int) {
	if len(s.tracks) <= track {
		return
	}
	s.tracks[track].active = !s.tracks[track].active
}

func (s *Sequencer) ToggleStep(track int, step int) {
	if len(s.tracks[track].steps) <= step {
		return
	}
	s.tracks[track].steps[step].active = !s.tracks[track].steps[step].active
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
