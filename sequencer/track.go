package sequencer

type Track struct {
	steps   []*Step
	device  int
	pulse   int
	channel uint8

	note     uint8
	velocity uint8
	active   bool
}

func (t Track) Pulse() int {
	return t.pulse
}

func (t Track) Steps() []*Step {
	return t.steps
}

func (t Track) ActiveStep() int {
	return t.pulse / (pulsesPerQuarterNote / stepsPerQuarterNote)
}

func (t *Track) incrPulse() {
	t.pulse++
	if t.pulse == pulsesPerQuarterNote*(stepsPerTrack/stepsPerQuarterNote) {
		t.resetPulse()
	}
	t.triggerStep()
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
		if i != t.ActiveStep() {
			step.reset()
			continue
		}
		step.trigger()
	}
}
