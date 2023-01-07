package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	trackTextActiveColor        = lipgloss.Color("232")
	trackTextInactiveColor      = lipgloss.Color("240")
	trackActiveColor            = lipgloss.Color("250")
	trackActiveStepTriggerColor = lipgloss.Color("255")
	trackStepTriggerColor       = lipgloss.Color("238")
	recModeColor                = lipgloss.Color("124")

	statusBarStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(0, 1, 0, 0).
			Bold(true)

	trackActiveStyle = statusBarStyle.Copy().
				Foreground(trackTextActiveColor).
				Background(trackActiveColor)
	trackActiveCurrentStepActiveStyle = statusBarStyle.Copy().
						Foreground(trackTextActiveColor).
						Background(trackActiveStepTriggerColor)
	trackCurrentStepActiveStyle = statusBarStyle.Copy().
					Background(trackStepTriggerColor)
	trackInactive = statusBarStyle.Copy().
			Foreground(trackTextInactiveColor)

	statusModeStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(recModeColor)

	tempoStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(recModeColor)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	logoStyle = lipgloss.NewStyle().Inherit(statusBarStyle)
)

func (m mainModel) renderStatus() string {
	w := lipgloss.Width

	statusTrack := m.renderStatusTracks()
	statusMode := statusModeStyle.Render("REC")
	logo := logoStyle.Render("SEKTRON")
	text := "Playing"
	statusVal := statusText.Copy().
		PaddingRight((m.width/stepsPerLine-2)*stepsPerLine - w(statusMode) - w(statusTrack) - w(logo) + w(text)).
		Render(text)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		statusTrack,
		statusMode,
		statusVal,
		logo,
	)
}

func (m mainModel) renderStatusTracks() string {
	var tracks []string
	for i, track := range m.seq.Tracks() {
		text := fmt.Sprintf("T%d", i+1)
		if i == m.activeTrack && track.IsCurrentStepActive() {
			tracks = append(tracks, trackActiveCurrentStepActiveStyle.Render(text))
		} else if m.seq.IsPlaying() && track.IsCurrentStepActive() {
			tracks = append(tracks, trackCurrentStepActiveStyle.Render(text))
		} else if i == m.activeTrack {
			tracks = append(tracks, trackActiveStyle.Render(text))
		} else if !track.IsActive() {
			tracks = append(tracks, trackInactive.Render(text))
		} else {
			tracks = append(tracks, statusBarStyle.Render(text))
		}

	}

	return lipgloss.JoinHorizontal(lipgloss.Center, tracks...)
}
