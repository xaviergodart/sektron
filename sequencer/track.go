package sequencer

type TrackInterface interface {
	Steps() []*step
	CurrentStep() int
	IsActive() bool
	IsCurrentStepActive() bool
}

type track struct {
	steps   []*step
	device  int
	pulse   int
	length  int
	channel uint8

	chord       []uint8
	velocity    uint8
	probability int
	trig        chan struct{}
	done        chan struct{}
	active      bool
}

func (t track) Steps() []*step {
	return t.steps
}

func (t track) CurrentStep() int {
	return t.pulse / pulsesPerStep
}

func (t track) IsActive() bool {
	return t.active
}

func (t track) IsCurrentStepActive() bool {
	if !t.active || len(t.steps) < t.CurrentStep() {
		return false
	}
	return t.steps[t.CurrentStep()].IsActive()
}

func (t *track) start() {
	t.trig = make(chan struct{})
	t.done = make(chan struct{})
	go func(track *track) {
		for {
			select {
			case <-track.trig:
				track.trigger()
			case <-track.done:
				return
			}
		}
	}(t)
}

func (t *track) tick() {
	t.trig <- struct{}{}
}

func (t track) stepForNextPulse() int {
	return (t.pulse + 1) % (pulsesPerStep * len(t.steps)) / pulsesPerStep
}

func (t track) isStepForNextPulseActive() bool {
	return t.steps[t.stepForNextPulse()].active
}

func (t *track) trigger() {
	for i, step := range t.steps {
		if t.active && step.isStartingPulse() {
			step.trigger()
		}
		// Avoid 2 steps to be triggered at the same time
		// when the first step overlaps
		if step.isEndingPulse() || (i != t.CurrentStep() && t.isStepForNextPulseActive()) {
			step.reset()
		}
	}
	t.pulse++
	if t.pulse == pulsesPerStep*len(t.steps) {
		t.pulse = 0
	}
}

func (t *track) reset() {
	t.pulse = 0
	for _, step := range t.steps {
		step.reset()
	}
}

func (t *track) close() {
	defer close(t.done)
	defer close(t.trig)
	t.done <- struct{}{}
}
