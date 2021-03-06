package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

	// List.

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Width(96 + 1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00b4d8"))

	tickCmd = func() tea.Cmd {
		return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	// Dialog.
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	typeButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#e76f51")).
			Padding(0, 3).
			MarginTop(1)

	familyButtonStyle = typeButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#90dbf4")).
				MarginRight(2)

	protocolButtonStyle = typeButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#2a9d8f")).
				MarginRight(2)

	width = 94
)
