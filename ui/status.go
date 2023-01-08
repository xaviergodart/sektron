package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Bold(true)

	trackActiveStyle = statusBarStyle.Copy().
				Foreground(secondaryTextColor).
				Background(activeColor)
	trackActiveCurrentStepActiveStyle = statusBarStyle.Copy().
						Foreground(secondaryTextColor).
						Background(currentColor)
	trackCurrentStepActiveStyle = statusBarStyle.Copy().
					Background(currentDimmedColor)
	trackInactive = statusBarStyle.Copy().
			Italic(true).
			Foreground(secondaryDimmedColor)

	statusPlayerStyle = statusBarStyle.Copy().
				Background(currentColor).
				Foreground(secondaryTextColor)
	statusPlayingStyle = statusPlayerStyle.Copy().
				Background(transportPlayColor)

	statusModeStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(transportRecColor)

	tempoStyle = statusBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(primaryColor)

	tempoTickStyle = tempoStyle.Copy().
			Foreground(primaryTextColor).
			Background(currentColor)

	statusTrackPage = lipgloss.NewStyle().
			MarginLeft(8)
	trackPage = statusBarStyle.Copy().
			Foreground(inactiveColor)
	trackPageActive = trackPage.Copy().
			Foreground(primaryColor)
	trackPageCurrent = trackPage.Copy().
				Foreground(activeColor)

	logoStyle = lipgloss.NewStyle().
			Italic(true).
			Inherit(statusBarStyle)
)

func (m mainModel) renderStatus() string {
	w := lipgloss.Width

	statusTrack := m.renderStatusTracks()
	statusMode := m.renderStatusMode()
	statusTempo := m.renderStatusTempo()
	statusPlayer := m.renderStatusPlayer()
	statusTrackPages := m.renderStatusTrackPages()

	statusBar := lipgloss.JoinHorizontal(lipgloss.Center,
		statusTempo,
		statusPlayer,
		statusMode,
		statusTrack,
		statusTrackPages,
	)

	logo := logoStyle.PaddingLeft((m.width/stepsPerLine-2)*stepsPerLine - w(statusBar) + 6).
		Render(sektron)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		statusBar,
		logo,
	)
}

func (m mainModel) renderStatusTracks() string {
	var tracks []string
	for i, track := range m.seq.Tracks() {
		text := fmt.Sprintf("T%d", i+1)
		if m.seq.IsPlaying() && i == m.activeTrack && track.IsCurrentStepActive() {
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

func (m mainModel) renderStatusTempo() string {
	text := fmt.Sprintf("⧗ %.1f", m.seq.Tempo())
	if m.isActiveTrackOnQuarterNote() {
		return tempoTickStyle.Render(text)
	}
	return tempoStyle.Render(text)
}

func (m mainModel) renderStatusPlayer() string {
	if m.seq.IsPlaying() {
		return statusPlayingStyle.Render("▶")
	}
	return statusPlayerStyle.Render("■")
}

func (m mainModel) renderStatusMode() string {
	text := "●"
	if m.mode == recMode {
		return statusModeStyle.Render(text)
	}
	return statusBarStyle.Render(text)
}

func (m mainModel) renderStatusTrackPages() string {
	pageNb := m.trackPagesNb()
	if pageNb <= 1 {
		return ""
	}
	text := "●"
	pages := make([]string, pageNb)
	for i := range pages {
		if m.isActiveTrackOnQuarterNote() && m.playingTrackPage() == i {
			pages[i] = trackPageCurrent.Render(text)
		} else if m.activeTrackPage == i {
			pages[i] = trackPageActive.Render(text)
		} else {
			pages[i] = trackPage.Render(text)
		}
	}
	return statusTrackPage.Render(
		lipgloss.JoinHorizontal(lipgloss.Left, pages...),
	)
}

func (m mainModel) isActiveTrackOnQuarterNote() bool {
	return m.seq.IsPlaying() && m.seq.Tracks()[0].CurrentStep()%4 == 0
}

func (m mainModel) playingTrackPage() int {
	return m.seq.Tracks()[m.activeTrack].CurrentStep() / stepsPerPage
}
