package midi

import (
	"errors"
	"log"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

type Server struct {
	devices   midi.OutPorts
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan midi.Message
	started   bool
}

func NewServer() (*Server, error) {
	devices := midi.GetOutPorts()
	if len(devices) == 0 {
		return nil, errors.New("no midi drivers")
	}
	return &Server{
		devices: devices,
		started: false,
	}, nil
}

func (s *Server) Start() error {
	var wg sync.WaitGroup
	wg.Add(len(s.devices))
	s.done = make(chan struct{})
	for i, device := range s.devices {
		s.outputs = append(s.outputs, make(chan midi.Message, 1000))
		go func(device drivers.Out, done <-chan struct{}, output <-chan midi.Message) {
			defer wg.Done()
			send, err := midi.SendTo(device)
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case <-done:
					return
				case note := <-output:
					send(note)
				}
			}
		}(device, s.done, s.outputs[i])
	}
	s.waitGroup = &wg
	s.started = true
	return nil
}

func (s *Server) NoteOn(device int, channel uint8, note uint8, velocity uint8) {
	// TODO: check if output exists. Handle hot plug?
	s.outputs[device] <- midi.NoteOn(channel, note, velocity)
}

func (s *Server) NoteOff(device int, channel uint8, note uint8) {
	// TODO: check if output exists. Handle hot plug?
	s.outputs[device] <- midi.NoteOff(channel, note)
}

func (s *Server) Close() {
	defer midi.CloseDriver()
	for range s.devices {
		s.done <- struct{}{}
	}
	s.waitGroup.Wait()
}
