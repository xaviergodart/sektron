// Package filesystem provides interfaces and serializable structures that
// allows saving/loading sequencer state to/from json files.
package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

const (
	maxPatterns = 64
)

// Bank holds a slice of patterns in memory
type Bank struct {
	Patterns []Pattern `json:"patterns"`
	Active   int       `json:"active"`
	filename string
}

// Pattern represents a sequencer state that is json serializable.
type Pattern struct {
	Tracks []Track `json:"tracks"`
	Tempo  float64 `json:"tempo"`
}

// IsFree returns true if the pattern is not used, false otherwise.
func (p Pattern) IsFree() bool {
	return p.Tracks == nil
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

// NewBank creates and loads a new bank from a given file.
func NewBank(filename string) Bank {
	bank := Bank{
		filename: filename,
		Patterns: make([]Pattern, maxPatterns),
	}
	bank.Load(filename)
	return bank
}

// Save serializes the Bank and writes it to a file.
func (b *Bank) Save() {
	content, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(b.filename, content, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Load reads a json and unmarshal its content to the Bank..
func (b *Bank) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	err = json.Unmarshal(content, b)
	if err != nil {
		log.Fatal(err)
	}
}
