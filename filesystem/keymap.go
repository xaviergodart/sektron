package filesystem

// KeyMap represents a keyboard mapping loaded from a json file.
type KeyMap struct {
	Tracks       [10]string `json:"tracks"`
	TracksToggle [10]string `json:"tracks_toggle"`
	Steps        [16]string `json:"steps"`
	StepsToggle  [16]string `json:"steps_toggle"`
	ParamMode    string     `json:"param_mode"`
	PatternMode  string     `json:"pattern_mode"`
	AddTrack     string     `json:"add_track"`
	RemoveTrack  string     `json:"remove_track"`
	AddStep      string     `json:"add_step"`
	RemoveStep   string     `json:"remove_step"`
	PageUp       string     `json:"page_up"`
	PageDown     string     `json:"page_down"`
	TempoUp      string     `json:"tempo_up"`
	TempoDown    string     `json:"tempo_down"`
	AddParam     string     `json:"add_param"`
	RemoveParam  string     `json:"remove_param"`
	Validate     string     `json:"validate"`
	Left         string     `json:"left"`
	Right        string     `json:"right"`
	Up           string     `json:"up"`
	Down         string     `json:"down"`
	Help         string     `json:"help"`
	Quit         string     `json:"quit"`
}

// NewDefaultAzertyKeyMap returns a new default KeyMap for azerty keyboards.
func NewDefaultAzertyKeyMap() KeyMap {
	return KeyMap{
		Tracks:       [10]string{"&", "é", "\"", "'", "(", "-", "è", "_", "ç", "à"},
		TracksToggle: [10]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		Steps:        [16]string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"},
		StepsToggle:  [16]string{"A", "Z", "E", "R", "T", "Y", "U", "I", "Q", "S", "D", "F", "G", "H", "J", "K"},
		ParamMode:    "tab",
		PatternMode:  "²",
		AddTrack:     "=",
		RemoveTrack:  ")",
		AddStep:      "+",
		RemoveStep:   "°",
		PageUp:       "p",
		PageDown:     "m",
		TempoUp:      "ctrl+up",
		TempoDown:    "ctrl+down",
		AddParam:     "shift+up",
		RemoveParam:  "shift+down",
		Validate:     "enter",
		Left:         "left",
		Right:        "right",
		Up:           "up",
		Down:         "down",
		Help:         "?",
	}
}

// NewDefaultAzertyMacKeyMap returns a new default KeyMap for azerty mac
// keyboards.
func NewDefaultAzertyMacKeyMap() KeyMap {
	return KeyMap{
		Tracks:       [10]string{"&", "é", "\"", "'", "(", "-", "è", "_", "ç", "à"},
		TracksToggle: [10]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		Steps:        [16]string{"a", "z", "e", "r", "t", "y", "u", "i", "q", "s", "d", "f", "g", "h", "j", "k"},
		StepsToggle:  [16]string{"A", "Z", "E", "R", "T", "Y", "U", "I", "Q", "S", "D", "F", "G", "H", "J", "K"},
		ParamMode:    "tab",
		PatternMode:  "@",
		AddTrack:     "-",
		RemoveTrack:  ")",
		AddStep:      "_",
		RemoveStep:   "°",
		PageUp:       "p",
		PageDown:     "m",
		TempoUp:      "shift+up",
		TempoDown:    "shift+down",
		AddParam:     "ctrl+up",
		RemoveParam:  "ctrl+down",
		Validate:     "enter",
		Left:         "left",
		Right:        "right",
		Up:           "up",
		Down:         "down",
		Help:         "?",
	}
}

// NewDefaultQwertyKeyMap returns a new default KeyMap for qwerty keyboards.
func NewDefaultQwertyKeyMap() KeyMap {
	return KeyMap{
		Tracks:       [10]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		TracksToggle: [10]string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")"},
		Steps:        [16]string{"q", "w", "e", "r", "t", "y", "u", "i", "a", "s", "d", "f", "g", "h", "j", "k"},
		StepsToggle:  [16]string{"Q", "W", "E", "R", "T", "Y", "U", "I", "A", "S", "D", "F", "G", "H", "J", "K"},
		ParamMode:    "tab",
		PatternMode:  "`",
		AddTrack:     "=",
		RemoveTrack:  "-",
		AddStep:      "+",
		RemoveStep:   "_",
		PageUp:       "p",
		PageDown:     ";",
		TempoUp:      "shift+up",
		TempoDown:    "shift+down",
		AddParam:     "ctrl+up",
		RemoveParam:  "ctrl+down",
		Validate:     "enter",
		Left:         "left",
		Right:        "right",
		Up:           "up",
		Down:         "down",
		Help:         "?",
	}
}

// NewDefaultQwertyMacKeyMap returns a new default KeyMap for qwerty mac
// keyboards.
func NewDefaultQwertyMacKeyMap() KeyMap {
	return KeyMap{
		Tracks:       [10]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		TracksToggle: [10]string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")"},
		Steps:        [16]string{"q", "w", "e", "r", "t", "y", "u", "i", "a", "s", "d", "f", "g", "h", "j", "k"},
		StepsToggle:  [16]string{"Q", "W", "E", "R", "T", "Y", "U", "I", "A", "S", "D", "F", "G", "H", "J", "K"},
		ParamMode:    "tab",
		PatternMode:  "§",
		AddTrack:     "=",
		RemoveTrack:  "-",
		AddStep:      "+",
		RemoveStep:   "_",
		PageUp:       "p",
		PageDown:     ";",
		TempoUp:      "shift+up",
		TempoDown:    "shift+down",
		AddParam:     "ctrl+up",
		RemoveParam:  "ctrl+down",
		Validate:     "enter",
		Left:         "left",
		Right:        "right",
		Up:           "up",
		Down:         "down",
		Help:         "?",
	}
}
