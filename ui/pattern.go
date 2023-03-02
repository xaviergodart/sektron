package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

const (
	patternWidth    = 15
	patternHeight   = patternWidth / 2
	patternsPerPage = 16
	patternsPerLine = 8
	patternPages    = 4
)

var (
	patternStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2).
			Bold(true).
			BorderStyle(lipgloss.HiddenBorder())

	patternCurrentStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Padding(0, 1, 0, 2).
				Bold(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(tertiaryColor)
)

func (m mainModel) renderPatterns() string {
	pages := make([][]string, patternPages)

	for i := range pages {
		pages[i] = make([]string, patternsPerPage)
		for j := 0; j < patternsPerPage; j++ {
			pages[i][j] = m.renderPattern(j + (i * patternsPerPage))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			pages[m.activePatternPage][:patternsPerLine]...,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			pages[m.activePatternPage][len(pages[m.activePatternPage])-patternsPerLine:]...,
		),
		m.renderChain(),
	)
}

func (m mainModel) renderPattern(pattern int) string {
	width, height := m.stepSize()
	number := strconv.Itoa(pattern + 1)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		number,
		toASCIIFont(fmt.Sprintf("P%s", number)),
	)
	if pattern == m.seq.ActivePattern() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(tertiaryColor).Render(content),
			lipgloss.WithWhitespaceBackground(tertiaryColor),
		))
	} else if !m.seq.Patterns()[pattern].IsFree() {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(activeColor).Render(content),
			lipgloss.WithWhitespaceBackground(activeColor),
		))
	} else {
		return stepStyle.Render(lipgloss.Place(
			width,
			height,
			lipgloss.Left,
			lipgloss.Top,
			textStyle.Background(inactiveColor).Render(number),
			lipgloss.WithWhitespaceBackground(inactiveColor),
		))
	}
}

func (m mainModel) renderChain() string {
	var patterns []string

	for i, pattern := range m.seq.FullChain() {
		if i == 0 {
			patterns = append(patterns, patternCurrentStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(fmt.Sprintf("P%d", pattern+1)),
					"",
					" current ",
				),
			))
		} else {
			patterns = append(patterns, patternStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					toASCIIFont(fmt.Sprintf("P%d", pattern+1)),
					"",
					fmt.Sprintf("next %d", i),
				),
			))
		}
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				patterns...,
			),
		)
}
