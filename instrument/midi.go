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
	midiChanBufferSize = 1024
)

type midiInstrument struct {
	devices   midi.OutPorts
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan midi.Message
	started   bool
}

func NewMidi() (Instrument, error) {
	devices := midi.GetOutPorts()
	if len(devices) == 0 {
		return nil, errors.New("no midi drivers")
	}
	instrument := &midiInstrument{
		devices: devices,
		started: false,
	}
	instrument.start()
	return instrument, nil
}

func (m *midiInstrument) start() {
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
	m.started = true
}

func (m *midiInstrument) NoteOn(device int, channel uint8, note uint8, velocity uint8) {
	m.outputs[device] <- midi.NoteOn(channel, note, velocity)
}

func (m *midiInstrument) NoteOff(device int, channel uint8, note uint8) {
	m.outputs[device] <- midi.NoteOff(channel, note)
}

func (m *midiInstrument) SendClock(devices []int) {
	for _, device := range devices {
		m.outputs[device] <- midi.TimingClock()
	}
}

func (m *midiInstrument) Close() {
	defer midi.CloseDriver()
	for range m.devices {
		m.done <- struct{}{}
	}
	m.waitGroup.Wait()
}
