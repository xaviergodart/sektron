package ui

import "github.com/charmbracelet/lipgloss"

func (u UI) ViewSequencer() string {
	var tracks []string
	for _, track := range u.seq.Tracks() {
		tracks = append(tracks, u.ViewTrack(track))
	}
	return lipgloss.JoinVertical(lipgloss.Left, tracks...)
}
