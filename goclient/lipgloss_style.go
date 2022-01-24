package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	url     = lipgloss.NewStyle().Foreground(special).Render

	highlight       = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	tab = lipgloss.NewStyle().
		Border(activeTabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	// Status Bar.

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

	clientIDKeyStyle = statusNugget.Copy().
				Background(lipgloss.Color("#A550DF")).
				Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	clientIDStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)
