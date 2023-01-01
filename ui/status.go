package ui

import "github.com/charmbracelet/lipgloss"

var (
	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	trackStatusStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

func (m mainModel) renderStatus() string {
	w := lipgloss.Width

	statusKey := statusStyle.Render("REC")
	statusTrack := trackStatusStyle.Render("TRACK 1")
	statusVal := statusText.Copy().
		Width(m.width - w(statusKey) - w(statusTrack)).
		Render("Playing")

	return lipgloss.JoinHorizontal(lipgloss.Top,
		statusTrack,
		statusKey,
		statusVal,
	)
}
