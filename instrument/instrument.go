// Package instrument provides ways to interact with music/audio devices and
// softwares.
//
// Right now, it only provides a midi instrument, that sends midi messages
// to capable devices. So, the main interface is closely linked to the midi
// protocol.
// It may evolves a lot if/when we start including other instruments.
package instrument

// Instrument provides a way to interct with an music instrument.
type Instrument interface {
	NoteOn(device int, channel uint8, note uint8, velocity uint8)
	NoteOff(device int, channel uint8, note uint8)
	Note(note uint8) string
	SendClock(devices []int)
	Close()
}
