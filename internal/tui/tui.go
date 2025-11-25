package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	width                   = 110
	columnWidth             = 30
	innerMargin             = 60
	columnKeyID             = "id"
	columnKeyReference      = "reference"
	columnKeyCustomer       = "customer"
	columnKeyAmount         = "amount"
	columnKeyStatus         = "status"
	columnKeyCreatedAt      = "created"
	columnKeyFirstName      = "firstname"
	columnKeyLastName       = "lastname"
	columnKeyEmail          = "email"
	columnKeyPhone          = "phone"
	columnKeyAddress        = "address"
	columnKeyCity           = "city"
	columnKeyCountry        = "country"
	columnKeyZip            = "zip"
	columnKeyNotes          = "notes"
	columnKeyCategory       = "category"
	columnKeyName           = "name"
	columnKeySku            = "sku"
	columnKeyPrice          = "price"
	columnKeyProductStatus  = "productstatus"
	columnKeyVoucherType    = "vouchertype"
	columnKeyVoucherCode    = "vouchercode"
	columnKeyValid          = "valid"
	columnKeyRequiredAmount = "requiredamount"
	columnKeyDiscount       = "discount"
	columnKeyExpires        = "expires"
)

var orderStatus = [3]string{"Pending", "Cancelled", "Paid"}
var productStatus = [4]string{"Instock", "Low Stock", "Out of Stock", "Preorder"}

type fn func(int) string

var (
	//t  table.Model
	vp viewport.Model

	normal = lipgloss.Color("#EEEEEE")
	base   = lipgloss.NewStyle().Foreground(normal)
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#585858"}

	background = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#CCCCCC"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#61D4C6"}

	divStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 2).
			BorderForeground(lipgloss.Color("69"))

	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			Align(lipgloss.Center)
	dialogTitleStyle = lipgloss.NewStyle().
				Inherit(titleBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
				Padding(0, 1)

	modalBodyStyle = lipgloss.NewStyle().
			Inherit(titleBarStyle).
			Foreground(lipgloss.Color("#FF5F87")).
			Background(lipgloss.Color("#000000")).
			Padding(0, 1)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

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

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

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

type model struct {
	cursor   int
	sections []string
	rowIndex int
}

func initialModel() model {
	return model{
		sections: []string{"Orders", "Customers", "Addresses", "Products", "Vouchers"},
		rowIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}
	//physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

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

		form := lipgloss.JoinVertical(lipgloss.Left, divStyle.Render("HELLO\nNew form: dsfasdfasdf\nType here: sadfasdfasdf\n"))
		body := lipgloss.JoinVertical(lipgloss.Left, bar, form)

		doc.WriteString(lipgloss.NewStyle().Width(width).Render(body) + "\n\n")
	}

	return docStyle.Render(doc.String())
}

func Run() {
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
