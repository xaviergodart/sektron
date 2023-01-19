package sequencer

import (
	"fmt"
	"sektron/midi"
)

type Parametrable interface {
	Chord() []uint8
	Length() int
	Velocity() uint8
	Probability() int
	SetChord(chord []uint8)
	SetLength(length int)
	SetVelocity(velocity uint8)
	SetProbability(probability int)
}

func ChordString(chord []uint8) string {
	return midi.Note(chord[0])
}

func LengthString(length int) string {
	switch length {
	case pulsesPerStep / 2:
		return "1/32"
	case pulsesPerStep:
		return "1/16"
	case pulsesPerStep * stepsPerQuarterNote:
		return "1/4"
	case pulsesPerStep * stepsPerQuarterNote * 3:
		return "3/4"
	default:
		return fmt.Sprintf("%.1f", float64(length)/float64(pulsesPerStep))
	}
}
