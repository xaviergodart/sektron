package sequencer

import "time"

const (
	pulsesPerStep       int = 6
	stepsPerQuarterNote int = 4
)

type Clock struct {
	ticker *time.Ticker
	update chan float64
}

func NewClock(tempo float64, tick func()) *Clock {
	clock := &Clock{
		ticker: time.NewTicker(NewClockInterval(tempo)),
		update: make(chan float64, 128),
	}
	go func(clock *Clock) {
		for {
			select {
			case <-clock.ticker.C:
				tick()
			case newTempo := <-clock.update:
				clock.ticker.Stop()
				clock.ticker = time.NewTicker(NewClockInterval(newTempo))
			}
		}
	}(clock)
	return clock
}

func NewClockInterval(tempo float64) time.Duration {
	return time.Duration(1000000*60/(tempo*float64(pulsesPerStep*stepsPerQuarterNote))) * time.Microsecond
}
