package midi

import (
	"errors"
	"log"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	// Weirdly, if you send clock to 2 midi devices at the same time
	// with portmidi, it crashes
	//_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

const (
	midiChanBufferSize = 1024
)

type MidiInterface interface {
	NoteOn(device int, channel uint8, note uint8, velocity uint8)
	NoteOff(device int, channel uint8, note uint8)
	SendClock(devices []int)
}

type Midi struct {
	devices   midi.OutPorts
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan midi.Message
	started   bool
}

func New() (*Midi, error) {
	devices := midi.GetOutPorts()
	if len(devices) == 0 {
		return nil, errors.New("no midi drivers")
	}
	instrument := &Midi{
		devices: devices,
		started: false,
	}
	instrument.start()
	return instrument, nil
}

func (m *Midi) start() error {
	var wg sync.WaitGroup
	wg.Add(len(m.devices))
	m.done = make(chan struct{})
	for i, device := range m.devices {
		m.outputs = append(m.outputs, make(chan midi.Message, midiChanBufferSize))
		go func(device drivers.Out, done <-chan struct{}, output <-chan midi.Message) {
			defer wg.Done()
			send, err := midi.SendTo(device)
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case <-done:
					// TODO: purge output buffer
					return
				case msg := <-output:
					send(msg)
				}
			}
		}(device, m.done, m.outputs[i])
	}
	m.waitGroup = &wg
	m.started = true
	return nil
}

func (m *Midi) NoteOn(device int, channel uint8, note uint8, velocity uint8) {
	m.outputs[device] <- midi.NoteOn(channel, note, velocity)
}

func (m *Midi) NoteOff(device int, channel uint8, note uint8) {
	m.outputs[device] <- midi.NoteOff(channel, note)
}

func (m *Midi) SendClock(devices []int) {
	for _, device := range devices {
		m.outputs[device] <- midi.TimingClock()
	}
}

func (m *Midi) Close() {
	defer midi.CloseDriver()
	for range m.devices {
		m.done <- struct{}{}
	}
	m.waitGroup.Wait()
}
