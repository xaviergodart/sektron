package sequencer

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/gomidi/midi/v2"
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
	tracks    []track
	tempo     float64
	pulse     int
	isPlaying bool
}

func New(sendMidi func(msg midi.Message) error) Sequencer {
	var steps []step
	for i := 1; i <= stepsPerTrack; i++ {
		steps = append(steps, step{
			note:      defaultNote,
			triggered: false,
		})
	}

	var tracks []track
	for i := 1; i <= 1; i++ {
		tracks = append(tracks, track{
			steps:    steps,
			sendMidi: sendMidi,
		})
	}
	return Sequencer{
		tracks:    tracks,
		tempo:     defaultTempo,
		pulse:     0.0,
		isPlaying: false,
	}
}

func (s *Sequencer) TogglePlay() {
	s.isPlaying = !s.isPlaying
	if !s.isPlaying {
		s.pulse = 0.0
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
	s.pulse++
	if s.pulse == pulsesPerQuarterNote*(stepsPerTrack/stepsPerQuarterNote) {
		s.pulse = 0.0
	}
}
