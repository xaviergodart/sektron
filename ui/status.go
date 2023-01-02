package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	recModeColor = lipgloss.Color("9")

	statusBarStyle = lipgloss.NewStyle().
		//BorderForeground(focusColor).
		Padding(1, 1).
		Bold(true)

	trackStatusStyle = statusBarStyle.Copy().
				Foreground(lipgloss.Color(primaryTextColor)).
				Padding(1, 1).
				Background(primaryColor)

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
	var tracks []string
	for i := range m.seq.Tracks() {
		tracks = append(tracks, trackStatusStyle.Render(fmt.Sprintf("T%d", i+1)))
	}

	statusTrack := lipgloss.JoinHorizontal(lipgloss.Center, tracks...)
	statusMode := statusModeStyle.Render("REC")
	logo := logoStyle.Render("SEKTRON")
	statusVal := statusText.Copy().
		PaddingRight(m.width - w(statusMode) - w(statusTrack) - w(logo)).
		Render("Playing")

	return lipgloss.JoinHorizontal(lipgloss.Center,
		statusTrack,
		statusMode,
		statusVal,
		logo,
	)
}
