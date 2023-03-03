package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	logoBig = []string{
		"░▒█▀▀▀█░▒█▀▀▀░▒█░▄▀░▀▀█▀▀░▒█▀▀▄░▒█▀▀▀█░▒█▄░▒█",
		"░░▀▀▀▄▄░▒█▀▀▀░▒█▀▄░░░▒█░░░▒█▄▄▀░▒█░░▒█░▒█▒█▒█",
		"░▒█▄▄▄█░▒█▄▄▄░▒█░▒█░░▒█░░░▒█░▒█░▒█▄▄▄█░▒█░░▀█",
	}

	transportBarStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Bold(true)

	trackActiveStyle = transportBarStyle.Copy().
				Foreground(secondaryTextColor).
				Background(activeColor)
	trackActiveInactiveStyle = transportBarStyle.Copy().
					Italic(true).
					Foreground(secondaryTextColor).
					Background(activeColor)
	trackActiveCurrentStepActiveStyle = transportBarStyle.Copy().
						Foreground(secondaryTextColor).
						Background(currentColor)
	trackCurrentStepActiveStyle = transportBarStyle.Copy().
					Background(currentDimmedColor)
	trackInactiveStyle = transportBarStyle.Copy().
				Italic(true).
				Foreground(secondaryDimmedColor)

	transportPlayerStyle = transportBarStyle.Copy().
				Background(currentColor).
				Foreground(secondaryTextColor)
	transportPlayingStyle = transportPlayerStyle.Copy().
				Background(tertiaryColor)

	tempoStyle = transportBarStyle.Copy().
			Foreground(primaryTextColor).
			Background(primaryColor)

	tempoTickStyle = tempoStyle.Copy().
			Foreground(primaryTextColor).
			Background(currentColor)

	transportTrackPage = lipgloss.NewStyle()
	trackPage          = transportBarStyle.Copy().
				Foreground(inactiveColor)
	trackPageActive = trackPage.Copy().
			Foreground(primaryColor)
	trackPageCurrent = trackPage.Copy().
				Foreground(activeColor)

	logoStyle = lipgloss.NewStyle().Foreground(logoColor)
)

func (m mainModel) renderTransport() string {
	w := lipgloss.Width

	transportTrack := m.renderTransportTracks()
	transportTempo := m.renderTransportTempo()
	transportPlayer := m.renderTransportPlayer()
	transportSignature := trackActiveCurrentStepActiveStyle.Render(
		fmt.Sprintf("%d/%d", len(m.getActiveTrack().Steps()), m.trackPagesNb()*stepsPerPage),
	)
	transportPages := m.renderTransportPages()

	transportBar := lipgloss.JoinHorizontal(lipgloss.Center,
		transportTempo,
		transportPlayer,
		transportTrack,
		transportSignature,
		transportPages,
	)

	logo := logoStyle.PaddingLeft((m.width/stepsPerLine-2)*stepsPerLine - w(transportBar) - w(logoBig[0]) + 14).
		Render(lipgloss.JoinVertical(lipgloss.Left, logoBig...))

	return lipgloss.JoinHorizontal(lipgloss.Center,
		transportBar,
		logo,
	)
}

func (m mainModel) renderTransportTracks() string {
	var tracks []string
	for i, track := range m.seq.Tracks() {
		text := fmt.Sprintf("T%d", i+1)
		if m.seq.IsPlaying() && i == m.activeTrack && track.IsCurrentStepActive() {
			tracks = append(tracks, trackActiveCurrentStepActiveStyle.Render(text))
		} else if m.seq.IsPlaying() && track.IsCurrentStepActive() {
			tracks = append(tracks, trackCurrentStepActiveStyle.Render(text))
		} else if i == m.activeTrack && !track.IsActive() {
			tracks = append(tracks, trackActiveInactiveStyle.Render(text))
		} else if i == m.activeTrack {
			tracks = append(tracks, trackActiveStyle.Render(text))
		} else if !track.IsActive() {
			tracks = append(tracks, trackInactiveStyle.Render(text))
		} else {
			tracks = append(tracks, transportBarStyle.Render(text))
		}

	}

	return lipgloss.JoinHorizontal(lipgloss.Center, tracks...)
}

func (m mainModel) renderTransportTempo() string {
	text := fmt.Sprintf("⧗ %.1f", m.seq.Tempo())
	if m.isActiveTrackOnQuarterNote() {
		return tempoTickStyle.Render(text)
	}
	return tempoStyle.Render(text)
}

func (m mainModel) renderTransportPlayer() string {
	if m.seq.IsPlaying() {
		return transportPlayingStyle.Render("▶")
	}
	return transportPlayerStyle.Render("■")
}

func (m mainModel) renderTransportPages() string {
	if m.mode == patternSelectMode {
		return m.renderTransportPatternPages()
	}

	return m.renderTransportTrackPages()
}

func (m mainModel) renderTransportPatternPages() string {
	text := "●"
	pages := make([]string, patternPages)
	for i := range pages {
		if m.activePatternPage == i {
			pages[i] = trackPageActive.Render(text)
		} else {
			pages[i] = trackPage.Render(text)
		}
	}
	return transportTrackPage.Render(
		lipgloss.JoinHorizontal(lipgloss.Left, pages...),
	)
}

func (m mainModel) renderTransportTrackPages() string {
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
	return transportTrackPage.Render(
		lipgloss.JoinHorizontal(lipgloss.Left, pages...),
	)
}

func (m mainModel) isActiveTrackOnQuarterNote() bool {
	return m.seq.IsPlaying() && m.seq.Tracks()[0].CurrentStep()%4 == 0
}

func (m mainModel) playingTrackPage() int {
	return m.getActiveTrack().CurrentStep() / stepsPerPage
}
