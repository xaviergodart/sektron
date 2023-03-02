package sequencer

import (
	"sektron/filesystem"
	"sektron/midi"
)

// Patterns returns all patterns from the bank.
func (s sequencer) Patterns() []filesystem.Pattern {
	return s.bank.Patterns
}

// ActivePattern returns the active pattern.
func (s sequencer) ActivePattern() int {
	return s.bank.Active
}

// GetPattern returns a new Pattern object from the current sequencer state.
func (s *sequencer) Save() {
	var tracks []filesystem.Track
	shouldSave := false
	for _, t := range s.Tracks() {
		var steps []filesystem.Step
		controls := map[int]int16{}
		for k := range t.activeControls {
			controls[k] = t.controls[k].Value()
		}

		for _, s := range t.Steps() {
			if !shouldSave && s.active {
				shouldSave = true
			}
			stepControls := map[int]int16{}
			for k, c := range s.controls {
				stepControls[k] = c.Value()
			}
			steps = append(steps, filesystem.Step{
				Active:      s.active,
				Controls:    stepControls,
				Length:      s.length,
				Chord:       s.chord,
				Velocity:    s.velocity,
				Probability: s.probability,
				Offset:      s.offset,
			})
		}

		tracks = append(tracks, filesystem.Track{
			Steps:       steps,
			Device:      t.device, // TODO: should we store the name instead?
			Channel:     t.channel,
			Controls:    controls,
			Length:      t.length,
			Chord:       t.chord,
			Velocity:    t.velocity,
			Probability: t.probability,
		})
	}

	if !shouldSave {
		return
	}

	s.bank.Patterns[s.bank.Active] = filesystem.Pattern{
		Tempo:  s.Tempo(),
		Tracks: tracks,
	}

	s.bank.Save()
}

// FullChain returns the full pattern chain.
func (s *sequencer) FullChain() []int {
	return append([]int{s.bank.Active}, s.chain...)
}

// Chain adds a pattern to the chain.
func (s *sequencer) Chain(pattern int) {
	s.Save()
	s.chain = append(s.chain, pattern)
}

// LoadNext empties the pattern chain and add the given pattern first in chain.
func (s *sequencer) ChainNow(pattern int) {
	s.Save()
	s.chain = make([]int, 1)
	s.chain[0] = pattern
}

// LoadNextInChain loads the first pattern in chain.
func (s *sequencer) LoadNextInChain() {
	if len(s.chain) > 0 {
		var pattern int
		pattern, s.chain = s.chain[0], s.chain[1:]
		s.Load(pattern)
	}
}

// LoadPattern loads a new sequencer state from Pattern object.
func (s *sequencer) Load(pattern int) {
	// close existing tracks first
	for _, t := range s.tracks {
		t.reset()
		t.close()
	}
	s.tracks = []*track{}
	s.bank.Active = pattern

	if s.bank.Patterns[pattern].Tracks == nil {
		for i := 0; i < defaultTracks; i++ {
			s.AddTrack()
		}
		return
	}

	s.SetTempo(s.bank.Patterns[pattern].Tempo)

	for i, t := range s.bank.Patterns[pattern].Tracks {
		// Check if midi device exists or set the first one found.
		if len(s.midi.Devices()) < t.Device+1 {
			t.Device = 0
		}

		s.tracks = append(s.tracks, &track{
			midi:                  s.midi,
			steps:                 []*step{},
			chord:                 t.Chord,
			length:                t.Length,
			velocity:              t.Velocity,
			probability:           t.Probability,
			device:                t.Device,
			channel:               t.Channel,
			activeControls:        map[int]struct{}{},
			lastSentControlValues: map[int]int16{},
			active:                true,
		})

		s.tracks[i].controls = midi.NewControls(s.midi, s.tracks[i])

		for k, v := range t.Controls {
			s.tracks[i].controls[k].Set(v)
			s.tracks[i].activeControls[k] = struct{}{}
		}

		s.tracks[i].steps = []*step{}
		for j, stp := range t.Steps {
			s.tracks[i].steps = append(s.tracks[i].steps, &step{
				position:    j,
				midi:        s.midi,
				track:       s.tracks[i],
				active:      stp.Active,
				length:      stp.Length,
				chord:       stp.Chord,
				velocity:    stp.Velocity,
				probability: stp.Probability,
				offset:      stp.Offset,
				controls:    map[int]*midi.Control{},
			})

			for k, v := range stp.Controls {
				s.tracks[i].steps[j].SetControl(k, v)
			}
		}

		s.tracks[i].start()
	}
}
