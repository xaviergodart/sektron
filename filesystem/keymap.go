package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// KeyMap represents a keyboard mapping loaded from a json file.
type KeyMap struct {
	Tracks        [10]string `json:"tracks"`
	TracksToggle  [10]string `json:"tracks_toggle"`
	Steps         [8]string  `json:"steps"`
	StepsToggle   [8]string  `json:"steps_toggle"`
	Play          string     `json:"play"`
	ParamMode     string     `json:"param_mode"`
	PatternMode   string     `json:"pattern_mode"`
	AddTrack      string     `json:"add_track"`
	RemoveTrack   string     `json:"remove_track"`
	AddStep       string     `json:"add_step"`
	RemoveStep    string     `json:"remove_step"`
	PageUp        string     `json:"page_up"`
	PageDown      string     `json:"page_down"`
	TempoUp       string     `json:"tempo_up"`
	TempoFineUp   string     `json:"tempo_fine_up"`
	TempoDown     string     `json:"tempo_down"`
	TempoFineDown string     `json:"tempo_fine_down"`
	AddParam      string     `json:"add_param"`
	RemoveParam   string     `json:"remove_param"`
	Validate      string     `json:"validate"`
	Left          string     `json:"left"`
	Right         string     `json:"right"`
	Up            string     `json:"up"`
	Down          string     `json:"down"`
}

// Load reads a json and unmarshal its content to the KeyMap.
func (k *KeyMap) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	err = json.Unmarshal(content, k)
	if err != nil {
		log.Fatal(err)
	}
}
