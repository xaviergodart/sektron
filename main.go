package main

import (
    "fmt"
    "os"
    "strconv"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var (
    stepStyle = lipgloss.NewStyle().
            Width(6).
            Height(3).
            Margin(1).
            Bold(true).
            Foreground(lipgloss.Color("#FAFAFA")).
            Background(lipgloss.Color("#7D56F4"))
    stepCurrentStyle = lipgloss.NewStyle().
                Width(6).
                Height(3).
                Margin(1).
                Bold(true).
                Foreground(lipgloss.Color("#000000")).
                Background(lipgloss.Color("#FFFFFF"))
)

type ClockTickMsg time.Time

type step struct {
    note int
}

type track struct {
    steps       []step
    currentStep int
}

func (t track) View() string {
    var steps []string
    for i, _ := range t.steps {
        if i == t.currentStep {
            steps = append(steps, stepCurrentStyle.Render(strconv.Itoa(i+1)))
        } else {
            steps = append(steps, stepStyle.Render(strconv.Itoa(i+1)))
        }
    }
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        steps...,
    )
}

type model struct {
    tracks    []track
    tempo     int
    isPlaying bool
}

func (m *model) TogglePlay() {
    m.isPlaying = !m.isPlaying
}

func (m model) Clock() tea.Cmd {
    if !m.isPlaying {
        return nil
    }
    return tea.Every(time.Duration(60000/m.tempo)*time.Millisecond, func(t time.Time) tea.Msg {
        return ClockTickMsg(t)
    })
}

func (m *model) UpdateTracks() {
    for i, track := range m.tracks {
        m.tracks[i].currentStep = (track.currentStep + 1) % len(track.steps)
    }
}

func initialModel() model {
    var steps []step
    for i := 1; i <= 16; i++ {
        steps = append(steps, step{
            note: 0,
        })
    }

    var tracks []track
    for i := 1; i <= 4; i++ {
        tracks = append(tracks, track{
            steps:       steps,
            currentStep: 0,
        })
    }
    return model{
        tracks:    tracks,
        tempo:     220,
        isPlaying: false,
    }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case ClockTickMsg:
        m.UpdateTracks()
        return m, m.Clock()

    case tea.KeyMsg:
        switch msg.String() {

        case " ":
            m.TogglePlay()
            return m, m.Clock()

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    var ui string

    // Tracks
    for _, track := range m.tracks {
        ui = ui + track.View()
    }
    // Send the UI for rendering
    return ui
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
