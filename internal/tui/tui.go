package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	width              = 96
	columnWidth        = 30
	columnKeyID        = "id"
	columnKeyReference = "reference"
	columnKeyCustomer  = "customer"
	columnKeyAmount    = "amount"
	columnKeyStatus    = "status"
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

	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "╥",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "╨",
		InnerJunction:  "╫",

		InnerDivider: "│",
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
	cursor     int
	sections   []string
	selected   int
	tableModel table.Model
}

func initialModel() model {
	ordersTable := GetOrders()
	return model{
		sections:   []string{"Orders", "Customers", "Addresses", "Products", "Vouchers"},
		selected:   0,
		tableModel: ordersTable,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func BlankTable() table.Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 10),
	}
	rows := []table.Row{}
	t := table.New(columns).
		WithRows(rows).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true))

	return t
}

func (m *model) updateTableModel() {
	if m.cursor > 0 {
		m.tableModel = BlankTable()
	} else {
		m.tableModel = GetOrders()
	}
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
			//case "up":
			//	m.selected++
			//case "down":
			//	m.selected--
		}
	}

	return m, nil
}

func (m model) View() string {
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

		m.updateTableModel()
		// Tabs
		// get model.selected to determine activeTab
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
		//vp := viewport.New(96, 20)
		//vp.YPosition = 20
		//vp.SetContent(m.tableModel.View())
		doc.WriteString(m.tableModel.View())
	}

	return docStyle.Render(doc.String())
}

func GetOrders() table.Model {
	orders, err := models.GetOrdersByStatus(0)
	if err != nil {
		log.Printf("ORDERS ERROR %v", err)
	}

	columns := []table.Column{
		table.NewColumn(columnKeyReference, "Reference", 15).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#88f")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyCustomer, "Customer", 10),
		table.NewColumn(columnKeyAmount, "Amount", 15),
		table.NewColumn(columnKeyStatus, "Status", 10),
	}
	rows := []table.Row{}

	for _, order := range orders {
		amount := fmt.Sprintf("%.2f", order.Amount/100.00)
		status := strconv.Itoa(int(order.Status))
		customerId := strconv.Itoa(order.CustomerId)
		newRow := table.NewRow(table.RowData{
			columnKeyReference: order.ReferenceCode,
			columnKeyCustomer:  customerId,
			columnKeyAmount:    amount,
			columnKeyStatus:    status,
		})
		rows = append(rows, newRow)
	}

	// Start with the default key map and change it slightly, just for demoing
	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	t = table.New(columns).
		WithRows(rows).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		SelectableRows(true).
		Focused(true).
		Border(customBorder).
		WithKeyMap(keys).
		WithStaticFooter("Footer!").
		WithPageSize(5).
		WithSelectedText(" ", "✓").
		WithBaseStyle(
			lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#a38")).
				Foreground(lipgloss.Color("#a7a")).
				Align(lipgloss.Left),
		)

	return t
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

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("セバエラー")
		os.Exit(1)
	}
}
