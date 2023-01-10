package sequencer

type Track struct {
	steps   []*Step
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

func (t *Track) start() {
	t.trig = make(chan struct{})
	t.done = make(chan struct{})
	go func(track *Track) {
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

func (t *Track) tick() {
	t.trig <- struct{}{}
}

func (t Track) Steps() []*Step {
	return t.steps
}

func (t Track) CurrentStep() int {
	return t.pulse / pulsesPerStep
}

func (t Track) IsActive() bool {
	return t.active
}

func (t Track) IsCurrentStepActive() bool {
	if !t.active || len(t.steps) < t.CurrentStep() {
		return false
	}
	return t.steps[t.CurrentStep()].IsActive()
}

func (t Track) stepForNextPulse() int {
	return (t.pulse + 1) % (pulsesPerStep * len(t.steps)) / pulsesPerStep
}

func (t Track) isStepForNextPulseActive() bool {
	return t.steps[t.stepForNextPulse()].active
}

func (t *Track) trigger() {
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
		t.reset()
	}
}

func (t *Track) reset() {
	t.pulse = 0
}

func (t *Track) clear() {
	t.reset()
	for _, step := range t.steps {
		step.reset()
	}
}

func (t *Track) close() {
	defer close(t.done)
	defer close(t.trig)
	t.done <- struct{}{}
}
