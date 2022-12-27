package midi

import (
	"errors"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Server struct {
	devices   midi.OutPorts
	waitGroup *sync.WaitGroup
	done      chan struct{}
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
	for _, device := range s.devices {
		go func(device drivers.Out) {
			defer wg.Done()
			<-s.done
		}(device)
	}
	s.waitGroup = &wg
	return nil
}

func (s Server) Close() {
	defer midi.CloseDriver()
	for range s.devices {
		s.done <- struct{}{}
	}
	s.waitGroup.Wait()
}
