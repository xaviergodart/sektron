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
	active      bool
}

func (t *Track) Start() {
	t.trig = make(chan struct{})
	go func(track *Track) {
		for {
			<-track.trig
			track.trigger()
		}
	}(t)
}

func (t *Track) Pulse() {
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

func (t Track) stepForNextPulse() int {
	return (t.pulse + 1) % (pulsesPerStep * len(t.steps)) / pulsesPerStep
}

func (t Track) isStepForNextPulseActive() bool {
	return t.steps[t.stepForNextPulse()].active
}

func (t *Track) trigger() {
	if !t.active {
		return
	}
	for i, step := range t.steps {
		if step.isStartingPulse() {
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
	for _, step := range t.steps {
		step.reset()
	}
}
