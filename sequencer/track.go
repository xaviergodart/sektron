package sequencer

import "gitlab.com/gomidi/midi/v2"

type track struct {
	steps    []step
	sendMidi func(msg midi.Message) error
}
