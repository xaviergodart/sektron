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

	tempoColor         = lipgloss.Color("27")
	recModeColor       = lipgloss.Color("124")
	playingStatusColor = lipgloss.Color("34")
	stoppedStatusColor = lipgloss.Color("250")

	statusBarStyle = lipgloss.NewStyle().
			Padding(1, 2).
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
			Italic(true).
			Foreground(trackTextInactiveColor)

	statusPlayerStyle = statusBarStyle.Copy().
				Background(stoppedStatusColor).
				Foreground(trackTextActiveColor)
	statusPlayingStyle = statusPlayerStyle.Copy().
				Background(playingStatusColor)

	statusModeStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(recModeColor)

	tempoStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(tempoColor)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	logoStyle = lipgloss.NewStyle().
			Italic(true).
			Inherit(statusBarStyle)
)

func (m mainModel) renderStatus() string {
	w := lipgloss.Width

	statusTrack := m.renderStatusTracks()
	statusMode := statusModeStyle.Render("●")
	statusTempo := tempoStyle.Render(fmt.Sprintf("⧗ %.1f", m.seq.Tempo()))
	var statusPlayer string
	if m.seq.IsPlaying() {
		statusPlayer = statusPlayingStyle.Render("▶")
	} else {
		statusPlayer = statusPlayerStyle.Render("■")
	}
	logo := logoStyle.PaddingLeft((m.width/stepsPerLine-2)*stepsPerLine - w(statusMode) - w(statusTrack) + w(statusPlayer) - w(sektron)).
		Render(sektron)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		statusTempo,
		statusPlayer,
		statusMode,
		statusTrack,
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
