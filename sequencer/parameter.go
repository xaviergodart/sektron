package sequencer

import (
	"fmt"
	"sektron/midi"
	"strconv"
)

const (
	minChordNote   = 21
	maxChordNote   = 108
	minLength      = 2
	maxLength      = pulsesPerStep*maxSteps + pulsesPerStep // +1 step for the infinity mode (TODO: fix INF on multiple tracks)
	minVelocity    = 0
	maxVelocity    = 127
	minProbability = 0
	maxProbability = 100
	minChannel     = 0
	maxChannel     = 15
	minOffset      = 0
	maxOffset      = 5
)

// Parametrable should be implemented by both step and track.
// Represents common parameter methods between both elements.
type Parametrable interface {
	Chord() []uint8
	Length() int
	Velocity() uint8
	Probability() int
	SetChord(chord []uint8)
	SetLength(length int)
	SetVelocity(velocity uint8)
	SetProbability(probability int)
	ChordString() string
	LengthString() string
	VelocityString() string
	ProbabilityString() string
	midi.Controllable
}

func chordString(chord []uint8) string {
	return midi.Note(chord[0])
}

func lengthString(length int) string {
	switch length {
	case pulsesPerStep / 2:
		return "1/32"
	case pulsesPerStep:
		return "1/16"
	case pulsesPerStep * stepsPerQuarterNote / 2:
		return "1/8"
	case pulsesPerStep * stepsPerQuarterNote:
		return "1/4"
	case pulsesPerStep * stepsPerQuarterNote * 2:
		return "1/2"
	case maxLength:
		return "INF"
	default:
		return fmt.Sprintf("%.1f", float64(length)/float64(pulsesPerStep))
	}
}

func velocityString(velocity uint8) string {
	return strconv.Itoa(int(velocity))
}

func probabilityString(probability int) string {
	return fmt.Sprintf("%d%%", probability)
}
