package sequencer

type Track struct {
	steps []*Step
	pulse int
}

func (t Track) GetPulse() int {
	return t.pulse
}

func (t Track) GetSteps() []*Step {
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
	for i, step := range t.steps {
		if i != t.ActiveStep() {
			step.reset()
			continue
		}
		step.trigger()
	}
}
