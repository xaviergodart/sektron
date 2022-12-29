package sequencer

import (
	"math/rand"

	"sektron/midi"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultTempo         float64 = 120.0
	defaultNote          uint8   = 60
	defaultVelocity      uint8   = 100
	defaultProbability   int     = 100
	defaultDevice        int     = 0
	defaultStepsPerTrack int     = 16

	pulsesPerStep int = 6
)

type ClockTickMsg time.Time

type Sequencer struct {
	midi      *midi.Server
	tracks    []*Track
	tempo     float64
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
			length:      pulsesPerStep / 2,
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
				//offset: i * pulsesPerStep,
			})
		}
		track.steps = steps
		tracks = append(tracks, track)
	}
	return &Sequencer{
		midi:      midi,
		tracks:    tracks,
		tempo:     defaultTempo,
		isPlaying: false,
	}
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

func (s Sequencer) Clock() tea.Cmd {
	if !s.isPlaying {
		return nil
	}
	// midi clock: http://midi.teragonaudio.com/tech/midispec/clock.htm
	return tea.Tick(time.Duration(60000000/(s.tempo*float64(pulsesPerStep*4)))*time.Microsecond, func(t time.Time) tea.Msg {
		return ClockTickMsg(t)
	})
}

func (s *Sequencer) Pulse() {
	if !s.isPlaying {
		return
	}
	for _, track := range s.tracks {
		track.incrPulse()
	}
	s.midi.SendClock()
}

func (s *Sequencer) Reset() {
	for _, track := range s.tracks {
		track.reset()
	}
}
