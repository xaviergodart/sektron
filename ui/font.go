package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	asciiCharTable = map[string]string{
		" ": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				" ",
				" ",
				" ",
			}...,
		),
		".": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				" ",
				" ",
				"▄",
			}...,
		),
		"/": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"  █",
				" █ ",
				"█  ",
			}...,
		),
		"%": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█  █ ",
				"  █  ",
				" █  █",
			}...,
		),
		"A": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				"█▄█",
				"█ █",
			}...,
		),
		"B": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				"█▀▄",
				"█▄█",
			}...,
		),
		"C": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				"█  ",
				"█▄█",
			}...,
		),
		"D": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▄",
				"█ █",
				"█▄▀",
			}...,
		),
		"E": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▀",
				"█▀ ",
				"█▄▄",
			}...,
		),
		"F": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▀",
				"█▀ ",
				"█  ",
			}...,
		),
		"G": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▀",
				"█ ▄",
				"█▄█",
			}...,
		),
		"I": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▀█▀",
				" █ ",
				"▄█▄",
			}...,
		),
		"N": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▄",
				"█ █",
				"█ █",
			}...,
		),
		"S": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				"▀▀▄",
				"█▄█",
			}...,
		),
		"T": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▀█▀",
				" █ ",
				" █ ",
			}...,
		),
		"b": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█  ",
				"█▄▄",
				"█▄█",
			}...,
		),
		"0": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				"█ █",
				"█▄█",
			}...,
		),
		"1": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▄█ ",
				" █ ",
				"▄█▄",
			}...,
		),
		"2": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				" ▄▀",
				"█▄▄",
			}...,
		),
		"3": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀█",
				" ▀▄",
				"█▄█",
			}...,
		),
		"4": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█ █",
				"█▄█",
				"  █",
			}...,
		),
		"5": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"█▀▀",
				"▀▀▄",
				"▄▄▀",
			}...,
		),
		"6": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▄▀▄",
				"█▄ ",
				"▀▄▀",
			}...,
		),
		"7": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▀▀█",
				" █ ",
				"▐▌ ",
			}...,
		),
		"8": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▄▀▄",
				"▄▀▄",
				"▀▄▀",
			}...,
		),
		"9": lipgloss.JoinVertical(
			lipgloss.Left,
			[]string{
				"▄▀▄",
				"▀▄█",
				" ▄▀",
			}...,
		),
	}
)

func toASCIIFont(str string) string {
	var chars []string
	for _, c := range str {
		if char, ok := asciiCharTable[string(c)]; ok {
			chars = append(chars, char, asciiCharTable[" "])
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, chars...)
}
