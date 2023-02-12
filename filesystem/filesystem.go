// Package filesystem provides interfaces and serializable structures that
// allows saving/loading sequencer state to/from json files.
package filesystem

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"encoding/json"
)

const (
	patternsPath = "patterns"
)

// Patternable should be implemented by the sequencer struct for state
// manipulation through Pattern objects.
type Patternable interface {
	GetPattern() Pattern
	LoadPattern(Pattern)
}

// Pattern represents a sequencer state that is json serializable.
type Pattern struct {
	Tempo  float64 `json:"tempo"`
	Tracks []Track `json:"tracks"`
}

// Track represents a sequencer track state that is json serializable.
type Track struct {
	Steps       []Step        `json:"steps"`
	Device      int           `json:"device"`
	Channel     uint8         `json:"channel"`
	Controls    map[int]int16 `json:"controls"`
	Length      int           `json:"length"`
	Chord       []uint8       `json:"chord"`
	Velocity    uint8         `json:"velocity"`
	Probability int           `json:"probability"`
}

// Step represents a sequencer step state that is json serializable.
type Step struct {
	Active      bool          `json:"active"`
	Controls    map[int]int16 `json:"controls"`
	Length      *int          `json:"length"`
	Chord       *[]uint8      `json:"chord"`
	Velocity    *uint8        `json:"velocity"`
	Probability *int          `json:"probability"`
	Offset      int           `json:"offset"`
}

// Save gets a pattern object from a patternable object (sequencer), serializes
// it, and writes it to a file.
func Save(name string, item Patternable) {
	os.MkdirAll(patternsPath, 0755)
	filename := fmt.Sprintf("%s/%s.json", patternsPath, name)
	content, err := json.Marshal(item.GetPattern())
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Load reads a json and make a pattern object from it, then loads it into a
// patternable object (sequencer).
func Load(name string, item Patternable) {
	filename := fmt.Sprintf("%s/%s.json", patternsPath, name)
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	pattern := Pattern{}
	json.Unmarshal(content, &pattern)

	item.LoadPattern(pattern)
}
