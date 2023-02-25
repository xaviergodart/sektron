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
	)
}

func (m mainModel) renderPattern(pattern int) string {
	// TODO: display chain somewhere
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
