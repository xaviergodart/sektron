package sequencer

type track struct {
	steps []*step
	pulse int
}

func (t track) activeStep() int {
	return t.pulse / (pulsesPerQuarterNote / stepsPerQuarterNote)
}

func (t *track) incrPulse() {
	t.pulse++
	if t.pulse == pulsesPerQuarterNote*(stepsPerTrack/stepsPerQuarterNote) {
		t.resetPulse()
	}
	t.triggerStep()
}

func (t *track) resetPulse() {
	t.pulse = 0
	for _, step := range t.steps {
		step.reset()
	}
}

func (t *track) triggerStep() {
	for i, step := range t.steps {
		if i != t.activeStep() {
			step.reset()
			continue
		}
		step.trigger()
	}
}
