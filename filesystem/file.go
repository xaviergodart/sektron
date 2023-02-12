package filesystem

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"encoding/json"
)

const (
	patternsPath = "patterns"
)

type Savable interface {
	SavablePattern() Pattern
	LoadPattern(Pattern)
}

type Pattern struct {
	Tempo  float64 `json:"tempo"`
	Tracks []Track `json:"tracks"`
}

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

type Step struct {
	Active      bool          `json:"active"`
	Controls    map[int]int16 `json:"controls"`
	Length      *int          `json:"length"`
	Chord       *[]uint8      `json:"chord"`
	Velocity    *uint8        `json:"velocity"`
	Probability *int          `json:"probability"`
	Offset      int           `json:"offset"`
}

func Save(name string, item Savable) {
	os.MkdirAll(patternsPath, 0755)
	filename := fmt.Sprintf("%s/%s.json", patternsPath, name)
	content, err := json.Marshal(item.SavablePattern())
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Load(name string, item Savable) {
	filename := fmt.Sprintf("%s/%s.json", patternsPath, name)
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := ioutil.ReadAll(f)
	pattern := Pattern{}
	json.Unmarshal(content, &pattern)

	item.LoadPattern(pattern)
}
