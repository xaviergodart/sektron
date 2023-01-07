package sequencer

import "time"

const (
	pulsesPerStep       int = 6
	stepsPerQuarterNote int = 4
)

type Clock struct {
	ticker *time.Ticker
	update chan float64
	tempo  float64
}

func NewClock(tempo float64, tick func()) *Clock {
	clock := &Clock{
		ticker: time.NewTicker(NewClockInterval(tempo)),
		update: make(chan float64, 128),
		tempo:  tempo,
	}
	go func(clock *Clock) {
		for {
			select {
			case <-clock.ticker.C:
				tick()
			case newTempo := <-clock.update:
				clock.ticker.Stop()
				clock.ticker = time.NewTicker(NewClockInterval(newTempo))
				clock.tempo = newTempo
			}
		}
	}(clock)
	return clock
}

func NewClockInterval(tempo float64) time.Duration {
	// midi clock: http://midi.teragonaudio.com/tech/midispec/clock.htm
	return time.Duration(1000000*60/(tempo*float64(pulsesPerStep*stepsPerQuarterNote))) * time.Microsecond
}
