package sequencer

type Track struct {
	steps   []*Step
	device  int
	pulse   int
	length  int
	channel uint8

	note     uint8
	velocity uint8
	active   bool
}

func (t Track) Steps() []*Step {
	return t.steps
}

func (t Track) CurrentStep() int {
	return t.pulse / pulsesPerStep
}

func (t Track) stepForNextPulse() int {
	return (t.pulse + 1) % (pulsesPerStep * len(t.steps)) / pulsesPerStep
}

func (t Track) isStepForNextPulseActive() bool {
	return t.steps[t.stepForNextPulse()].active
}

func (t *Track) incrPulse() {
	t.triggerStep()
	t.pulse++
	if t.pulse == pulsesPerStep*len(t.steps) {
		t.resetPulse()
	}
}

func (t *Track) resetPulse() {
	t.pulse = 0
	for _, step := range t.steps {
		step.reset()
	}
}

func (t *Track) triggerStep() {
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
}
