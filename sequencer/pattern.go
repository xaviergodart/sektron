package sequencer

import (
	"sektron/filesystem"
	"sektron/midi"
)

// GetPattern returns a new Pattern object from the current sequencer state.
func (s sequencer) GetPattern() filesystem.Pattern {
	var tracks []filesystem.Track
	for _, t := range s.Tracks() {
		var steps []filesystem.Step
		controls := map[int]int16{}
		for k := range t.activeControls {
			controls[k] = t.controls[k].Value()
		}

		for _, s := range t.Steps() {
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
			Device:      t.device,
			Channel:     t.channel,
			Controls:    controls,
			Length:      t.length,
			Chord:       t.chord,
			Velocity:    t.velocity,
			Probability: t.probability,
		})
	}

	return filesystem.Pattern{
		Tempo:  s.Tempo(),
		Tracks: tracks,
	}
}

// LoadPattern loads a new sequencer state from Pattern object.
func (s *sequencer) LoadPattern(pattern filesystem.Pattern) {
	// close existing tracks first
	for _, t := range s.tracks {
		t.close()
	}
	s.tracks = []*track{}
	s.SetTempo(pattern.Tempo)

	for i, t := range pattern.Tracks {
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
