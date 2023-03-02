# Sektron

Sektron is a midi step sequencer, made with live performance in mind, that runs in the terminal (TUI). It's heavily inspired by [elektron machines](https://www.elektron.se).

**_Sektron is still under heavy development. Features are missing and it's probably unstable._**

![sektron screenshot](/docs/screenshot.png)

## Build

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
./bin/sektron
```

Hit `?` to see all keybindings. `esc` to quit.

![sektron gif](/docs/vhs.gif)


## Roadmap for 0.1

 - [x] Exhaustive midi messages type (CC etc...)
 - [x] Basic pattern management (and chaining)
 - [ ] Improve controls UX
 - [ ] Key mapping management and configuration

## Roadmap for later

 - Step polyphony
 - Keyboard mode
 - Live record
 - More random/generative features
 - Retrigs
 - LFOs
