package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	width       = 96
	columnWidth = 30
)

type fn func(int) string

var (
	t  table.Model
	vp viewport.Model

	normal = lipgloss.Color("#EEEEEE")
	base   = lipgloss.NewStyle().Foreground(normal)

	background = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#CCCCCC"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#61D4C6"}

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle).MarginLeft(1)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true).BorderForeground(highlight)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type model struct {
	cursor   int
	sections []string
	selected int
}

func initialModel() model {
	return model{
		sections: []string{"Orders", "Customers", "Addresses", "Products", "Vouchers"},
		selected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "left", "<":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.sections) - 1
			}
		case "right", "tab", ">":
			if m.cursor < len(m.sections)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "up":
			m.selected++
		case "down":
			m.selected--
		}
	}

	return m, nil
}

func (m model) View() string {
	doc := strings.Builder{}

	// Tabs
	// get model.selected to determine activeTab
	{
		w := lipgloss.Width
		leftStatus := statusStyle.Render("<<<<")
		rightStatus := statusStyle.Render(">>>>")
		statusVal := statusText.
			Width(width - w(leftStatus) - w(rightStatus) - 1).Render("SHOP DASHBOARD")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			statusVal,
			rightStatus,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar) + "\n\n")

		var nav [5]string
		for i, menuItem := range m.sections {
			if m.cursor == i {
				nav[i] = activeTab.Render(menuItem)
			} else {
				nav[i] = tab.Render(menuItem)
			}
		}
		row := lipgloss.JoinHorizontal(
			lipgloss.Top, nav[0], nav[1], nav[2], nav[3], nav[4],
		)

		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	{
		var activeContent fn
		switch m.cursor {
		case 0:
			activeContent = GetOrders
		default:
			activeContent = Blank
		}
		vp := viewport.New(96, 20)
		vp.YPosition = 20
		vp.SetContent(SetActiveContent(activeContent, m.selected))
		doc.WriteString(vp.View())
	}

	return docStyle.Render(doc.String())
}

func SetActiveContent(f fn, s int) string {
	return f(s)
}

func Blank(int) string {
	str := ""
	return str
}

func GetOrders(selected int) string {
	columns := []table.Column{
		{Title: "Reference", Width: 10},
		{Title: "Customer", Width: 10},
		{Title: "Amount", Width: 15},
		{Title: "Status", Width: 10},
	}

	orders, err := models.GetOrdersByStatus(0)
	if err != nil {
		log.Printf("ORDERS ERROR %v", err)
	}
	rows := []table.Row{}

	for i, order := range orders {
		amount := fmt.Sprintf("%.2f", order.Amount/100.00)
		status := strconv.Itoa(int(order.Status))
		customerId := strconv.Itoa(order.CustomerId)
		newRow := []string{order.ReferenceCode, customerId, amount, status}
		newRow = newRow.NewRow(newRow)
		if i == selected {
			newRow = newRow.Selected(true)
		}
		rows = append(rows, newRow)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t.View()
}

func PrintOrders() {
	doc := strings.Builder{}
	{
		w := lipgloss.Width
		leftStatus := statusStyle.Render("<<<<")
		rightStatus := statusStyle.Render(">>>>")
		statusVal := statusText.
			Width(width - w(leftStatus) - w(rightStatus) - 1).Render("SHOP DASHBOARD")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			statusVal,
			rightStatus,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar) + "\n\n")
	}
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Orders"),
			tab.Render("Customers"),
			tab.Render("Addresses"),
			tab.Render("Products"),
			tab.Render("Inventory"),
			tab.Render("Vouchers"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	{
		vp := viewport.New(96, 30)
		//vp.YPosition = 20
		columns := []table.Column{
			{Title: "Reference", Width: 10},
			{Title: "Customer", Width: 10},
			{Title: "Amount", Width: 15},
			{Title: "Status", Width: 10},
		}

		orders, err := models.GetOrdersByStatus(0)
		if err != nil {
			log.Printf("ORDERS ERROR %v", err)
		}
		rows := []table.Row{}

		for _, order := range orders {
			amount := fmt.Sprintf("%.2f", order.Amount/100.00)
			status := strconv.Itoa(int(order.Status))
			customerId := strconv.Itoa(order.CustomerId)
			newRow := []string{order.ReferenceCode, customerId, amount, status}
			rows = append(rows, newRow)
		}

		t = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(20),
		)
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)
		vp.SetContent(t.View())
		doc.WriteString(vp.View())
	}

	fmt.Println(docStyle.Render(doc.String()))
}

func Run() {
	//f, err := tea.LogToFile("debug.log", "debug")
	//if err != nil {
	//	fmt.Println("fatal: ", err)
	//	os.Exit(1)
	//}
	//defer f.Close()

	cfg, err := config.LoadConfig(".env")
	_ = models.ConnectDatabase(&cfg)
	if err != nil {
		log.Printf("ERROR LOADING CONFIG")
	}

	//PrintOrders()

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("セバエラー")
		os.Exit(1)
	}
}
