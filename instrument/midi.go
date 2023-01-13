package instrument

import (
	"errors"
	"log"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

const (
	// Each midi device can receive notes through a dedicated buffered chan.
	// 16 tracks with all steps activated sending notes to the same device
	// at high tempo can results to a lot of midi messages.
	midiBufferSize = 1024
)

// midiInstrument contains the midi instrument state. We use the gomidi package
// for communicating with available midi devices.
type midiInstrument struct {
	// devices holds all the midi devices outputs that are returned by gomidi.
	devices midi.OutPorts

	// Because we want to allow the usage of multiple midi devices at the same
	// time, we start a goroutine for each device that can receive note trigs.
	// The wait group is used when closing the instrument (waits for all
	// device goroutines to end).
	// The done chan is used to send the end signal to the goroutines.
	// The output chans receives actual midi messages for each instruments.
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan midi.Message
}

// NewMidi creates a new midi instrument. It retrieves the connected midi
// devices and starts a new goroutine for each of them.
func NewMidi() (Instrument, error) {
	devices := midi.GetOutPorts()
	if len(devices) == 0 {
		return nil, errors.New("no midi drivers")
	}
	instrument := &midiInstrument{
		devices: devices,
	}
	instrument.start()
	return instrument, nil
}

func (m *midiInstrument) start() {
	var wg sync.WaitGroup
	wg.Add(len(m.devices))
	m.done = make(chan struct{}, len(m.devices))
	for i, device := range m.devices {
		m.outputs = append(m.outputs, make(chan midi.Message, midiBufferSize))
		go func(device drivers.Out, done <-chan struct{}, output <-chan midi.Message) {
			defer wg.Done()
			send, err := midi.SendTo(device)
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case <-done:
					// Before terminating the goroutine, we drain all the
					// remaining messages, ensuring that all the note off
					// signals will be sent before exiting.
					for len(output) > 0 {
						send(<-output)
					}
					return
				case msg := <-output:
					err := send(msg)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}(device, m.done, m.outputs[i])
	}
	m.waitGroup = &wg
}

// NoteOn sends a Note On midi meessage to the given device.
func (m *midiInstrument) NoteOn(device int, channel uint8, note uint8, velocity uint8) {
	m.outputs[device] <- midi.NoteOn(channel, note, velocity)
}

// NoteOff sends a Note Off midi meessage to the given device.
func (m *midiInstrument) NoteOff(device int, channel uint8, note uint8) {
	m.outputs[device] <- midi.NoteOff(channel, note)
}

// SendClock sends a Clock midi meessage to given devices.
func (m *midiInstrument) SendClock(devices []int) {
	for _, device := range devices {
		m.outputs[device] <- midi.TimingClock()
	}
}

// Close terminates all the device goroutines gracefully.
func (m *midiInstrument) Close() {
	defer midi.CloseDriver()
	if m.waitGroup == nil {
		return
	}
	for range m.devices {
		m.done <- struct{}{}
	}
	m.waitGroup.Wait()
}