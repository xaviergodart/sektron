package sequencer

import "github.com/charmbracelet/lipgloss"

func (s Sequencer) View() string {
	var tracks []string
	for _, track := range s.tracks {
		tracks = append(tracks, track.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, tracks...)
}
