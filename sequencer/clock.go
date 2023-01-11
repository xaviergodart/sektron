package sequencer

import "time"

const (
	pulsesPerStep       int     = 6
	stepsPerQuarterNote int     = 4
	tempoMin            float64 = 1.0
	tempoMax            float64 = 300.0
	updateBufferSize    int     = 128
)

type clock struct {
	ticker *time.Ticker
	update chan float64
	tempo  float64
}

func (c *clock) setTempo(tempo float64) {
	if tempo > tempoMax || tempo < tempoMin {
		return
	}
	c.update <- tempo
}

func newClock(tempo float64, tick func()) *clock {
	c := &clock{
		ticker: time.NewTicker(newClockInterval(tempo)),
		update: make(chan float64, updateBufferSize),
		tempo:  tempo,
	}
	go func(c *clock) {
		for {
			select {
			case <-c.ticker.C:
				tick()
			case newTempo := <-c.update:
				c.ticker.Stop()
				c.ticker = time.NewTicker(newClockInterval(newTempo))
				c.tempo = newTempo
			}
		}
	}(c)
	return c
}

func newClockInterval(tempo float64) time.Duration {
	// midi clock: http://midi.teragonaudio.com/tech/midispec/clock.htm
	return time.Duration(1000000*60/(tempo*float64(pulsesPerStep*stepsPerQuarterNote))) * time.Microsecond
}
