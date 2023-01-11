package instrument

type Instrument interface {
	NoteOn(device int, channel uint8, note uint8, velocity uint8)
	NoteOff(device int, channel uint8, note uint8)
	SendClock(devices []int)
	Close()
}
