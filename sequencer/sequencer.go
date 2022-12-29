package sequencer

import (
	"sektron/midi"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultTempo         = 120.0
	defaultNote          = 60
	defaultVelocity      = 100
	defaultDevice        = 0
	defaultStepsPerTrack = 16

	pulsesPerStep = 24
)

type ClockTickMsg time.Time

type Sequencer struct {
	midi      *midi.Server
	tracks    []*Track
	tempo     float64
	isPlaying bool
}

func New(midi *midi.Server) *Sequencer {
	var tracks []*Track
	for i := 0; i <= 1; i++ {
		var steps []*Step
		track := &Track{
			pulse:    0,
			note:     defaultNote + uint8(i*12) + uint8(i*5),
			length:   pulsesPerStep * 2,
			velocity: defaultVelocity,
			device:   defaultDevice,
			channel:  uint8(i),
			active:   true,
		}
		for j := 0; j < defaultStepsPerTrack; j++ {
			steps = append(steps, &Step{
				midi:   midi,
				track:  track,
				active: j%4 == 0,
			})
		}
		track.steps = steps
		tracks = append(tracks, track)
	}
	return &Sequencer{
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
	return tea.Every(time.Duration(60000000/(s.tempo*pulsesPerStep*4))*time.Microsecond, func(t time.Time) tea.Msg {
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
}

func (s *Sequencer) Reset() {
	for _, track := range s.tracks {
		track.resetPulse()
	}
}
