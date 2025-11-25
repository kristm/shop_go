package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
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

	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	noStyle             = lipgloss.NewStyle()

	normal = lipgloss.Color("#EEEEEE")
	base   = lipgloss.NewStyle().Foreground(normal)
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#585858"}

	background = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#CCCCCC"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#61D4C6"}

	baseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFF7DB"))

	divStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 2).
			BorderForeground(lipgloss.Color("69"))

	divItem = baseStyle.Foreground(lipgloss.Color("#FFD046")).Render

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
	cursor     int
	inputs     []textinput.Model
	focusIndex int
	product    models.Product
	cursorMode cursor.Mode
}

func initialModel(product *models.Product) model {
	m := model{
		inputs:     make([]textinput.Model, 5),
		focusIndex: 0,
		cursorMode: cursor.CursorBlink,
		product:    *product,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32
		t.Width = 70
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = product.Name
			t.TextStyle = focusedStyle.Foreground(lipgloss.Color("#ffffff"))
			t.PromptStyle = focusedStyle
			t.Focus()
		case 1:
			t.Placeholder = product.Sku
		case 2:
			t.Placeholder = product.Description
		case 3:
			t.Placeholder = fmt.Sprintf("%.0f", product.Price)
		case 4:
			t.Placeholder = fmt.Sprintf("%d", product.CategoryId)
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var (
	//	cmds []tea.Cmd
	//)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}
	//physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	{
		w := lipgloss.Width
		leftStatus := statusStyle.Render("<<<<")
		rightStatus := statusStyle.Render(">>>>")
		statusVal := statusText.
			Width(width - w(leftStatus) - w(rightStatus) - 1).Render(fmt.Sprintf("Product ID: %d", m.product.Id))

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			statusVal,
			rightStatus,
		)

		div := divStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				divItem(fmt.Sprintf("%s%s", "Name:", m.inputs[0].View())),
				divItem(fmt.Sprintf("%s%s", "SKU:", m.inputs[1].View())),
				divItem(fmt.Sprintf("%s%s", "Description:", m.inputs[2].View())),
				divItem(fmt.Sprintf("%s%s", "Price:", m.inputs[3].View())),
				divItem(fmt.Sprintf("%s%s", "CategoryId:", m.inputs[4].View())),
			),
		)

		body := lipgloss.JoinVertical(lipgloss.Left, bar, div)

		doc.WriteString(lipgloss.NewStyle().Width(width).Render(body) + "\n\n")
	}

	return docStyle.Render(doc.String())
}

func Run(ref string) {
	cfg, err := config.LoadConfig(".env")
	_ = models.ConnectDatabase(&cfg)
	if err != nil {
		log.Printf("ERROR LOADING CONFIG")
	}

	product, _ := models.GetProductBySku(ref)

	if _, err := tea.NewProgram(initialModel(&product)).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
