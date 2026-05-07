package modes

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var (
	//t  table.Model
	vp            viewport.Model
	ColumnWidth   int
	PhysicalWidth int

	focusedStyle = lipgloss.NewStyle().Foreground(CYAN)
	blurredStyle = lipgloss.NewStyle().Foreground(DARKGRAY)
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	baseStyle = lipgloss.NewStyle().
			Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(CYAN).
			Inherit(baseStyle).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 1).
			BorderForeground(DEEPPINK)

	buttonBlurredStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				Bold(false).
				Padding(0, 1).
				BorderForeground(DARKGRAY)

	divStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 2).
			BorderForeground(CORNFLOWERBLUE)
	divBlurredStyle = divStyle.BorderForeground(GRAY)

	divItem        = baseStyle.Foreground(GOLD).Render
	divBlurredItem = baseStyle.Foreground(GRAY).Render

	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			Align(lipgloss.Center)
	dialogTitleStyle = lipgloss.NewStyle().
				Inherit(titleBarStyle).
				Foreground(WHITE).
				Background(DEEPPINK).
				Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#FFFDF5"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle).MarginLeft(1)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)
