package sequencer

import (
	"sektron/midi"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultTempo = 90.0
	defaultNote  = 60

	pulsesPerQuarterNote = 24
	stepsPerQuarterNote  = 4
	stepsPerTrack        = 16
)

type ClockTickMsg time.Time

type Sequencer struct {
	midi      *midi.Server
	tracks    []*Track
	tempo     float64
	isPlaying bool
}

func New(midi *midi.Server) *Sequencer {
	var steps []*Step
	for i := 1; i <= stepsPerTrack; i++ {
		steps = append(steps, &Step{
			midi: midi,
			note: defaultNote + uint8(i),
		})
	}

	var tracks []*Track
	for i := 1; i <= 1; i++ {
		tracks = append(tracks, &Track{
			steps: steps,
			pulse: 0,
		})
	}
	return &Sequencer{
		tracks:    tracks,
		tempo:     defaultTempo,
		isPlaying: false,
	}
}

func (s *Sequencer) GetTracks() []*Track {
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
	return tea.Every(time.Duration(60000000/(s.tempo*pulsesPerQuarterNote))*time.Microsecond, func(t time.Time) tea.Msg {
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
