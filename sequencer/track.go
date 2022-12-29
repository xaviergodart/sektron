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

func (t Track) relativePulseForStep(step int) int {
	return t.pulse - (step * pulsesPerStep)
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
		if i == t.CurrentStep() && t.relativePulseForStep(i) == 0 {
			step.trigger()
		}
		if t.relativePulseForStep(i) >= step.Length() {
			step.reset()
		}
	}
}
