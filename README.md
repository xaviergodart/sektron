# Sektron

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/xaviergodart/sektron) ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/xaviergodart/sektron/build.yml)

Sektron is a midi step sequencer, made with live performance in mind, that runs in the terminal (TUI).
Its main purpose is to mimic the fun and immediacity of hardware sequencers by being entirely controllable via keyboard.

It's heavily inspired by [elektron devices](https://www.elektron.se).

**_Sektron is still under heavy development. Features are missing and it's probably unstable._**

![sektron screenshot](/docs/screenshot.png)

### Features

 - Fully (and only) **controllable by keyboard**
 - **Customizable** keyboard mapping
 - Up to **10 midi tracks**, that can be attached to specific midi device and channel
 - Up to **128 steps per track**. The number of steps per track is independent, allowing complex polyrhythms
 - Parameters can be set per track or step (**parameter locking**)
 - Up to **64 patterns** can be loaded at the same time
 - **Pattern chaining**

And more to come. See [Planned features](https://github.com/xaviergodart/sektron#planned-features).

## Installation

[Download the last release](https://github.com/xaviergodart/sektron/releases) for your platform.

Then:
```sh
# Extract files
mkdir -p sektron && tar -zxvf sektron_VERSION_PLATFORM.tar.gz -C sektron

# Run sektron
./sektron
```

### Build it yourself

You'll need [go 1.18](https://go.dev/dl/) minimum.
Although you should be able to build it for either **linux**, **macOS** or **Windows**, it has only been tested on **linux**.

```sh
# Linux
make GOLANG_OS=linux build

# macOS
make GOLANG_OS=darwin build

# Windows
make GOLANG_OS=windows build
```

## Usage

```sh
# Run sektron
./sektron

# Display current version
./sektron --version
```

Hit `?` to see all keybindings. `esc` to quit.

![sektron gif](/docs/vhs.gif)

[Qsynth](https://qsynth.sourceforge.io/) is a nice companion app for testing Sektron.


### Keyboard mapping

Keys mapping is fully customizable. After running Sektron for the first time, a `config.json` is created.
You can edit all the keys inside it.

You can select one of the fex default keyboard layouts are available:
```sh
# QWERTY
./sektron --keyboard qwerty

# AZERTY
./sektron --keyboard azerty

# QWERTY MAC
./sektron --keyboard qwerty-mac

# AZERTY MAC
./sektron --keyboard azerty-mac
```

### Default keyboard mapping

The default key mapping looks like this:

![keyboard layout](/docs/keyboard-layout.png)

For qwerty keyboards, here's the default mapping:

 - `space` **play** or **stop**
 - `tab` **toggle parameter mode (track, record)**
 - `` ` `` **toggle pattern select mode**
 - `1` `2` `3` `4` `5` `6` `7` `8` `9` `0` **select track**
 - `!` `@` `#` `$` `%` `^` `&` `*` `(` `)` **toggle track**
 - `q` `w` `e` `r` `t` `y` `u` `i` `q` `s` `d` `f` `g` `h` `j` `k` **select step** or **switch to pattern**
 - `Q` `W` `E` `R` `T` `Y` `U` `I` `Q` `S` `D` `F` `G` `H` `J` `K` **toggle step** or **chain pattern**
 - `=` **add track**
 - `-` **remove track**
 - `+` **add step** to the selected track
 - `_` **remove step** form the selected track
 - `p` **page up** either steps or patterns if more than 16 items
 - `;` **page down** either steps or patterns if more than 16 items
 - `shift`+`up` **increase tempo**
 - `shift`+`down` **decrease tempo**
 - `ctrl`+`up` **add new midi control** to the selected track
 - `ctrl`+`down` **remove midi control**. It will remove the selected one
 - `enter` **validate selection**
 - `left` **move parameter selection to the left**
 - `right` **move parameter selection to the right**
 - `up` **increase selected parameter value**
 - `down` **decrease selected parameter value**
 - `?` **show help**
 - `escape` or `ctrl`+`c` **quit**

### Patterns management

Each time you start Sektron, a json file (default: `patterns.json`) containing 128 pattern slots is loaded.
For selecting a different file, use the `--patterns` flag:
```sh
./sektron --patterns my-patterns.json
```

Each time you change pattern or quit the program, the current pattern is saved to the file.


## Acknowledgments

Sektron uses a few awesome packages:
 - [gomidi/midi](https://gitlab.com/gomidi/midi) for all midi communication
 - [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) as the main TUI framework
 - [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) for making things beautiful

## Planned features

### v0.1

 - [x] Exhaustive midi messages type (CC etc...)
 - [x] Basic pattern management (and chaining)
 - [x] Key mapping management and configuration
 - [ ] Improve controls UX
 - [ ] Write documentation


### Considered

 - Step polyphony
 - Keyboard mode
 - Live record
 - More random/generative features
 - Retrigs
 - LFOs
