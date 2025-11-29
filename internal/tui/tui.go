package tui

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	width      = 80
	MAX_INDEX  = 5
	FORM_WIDTH = 60
	timeout    = time.Second * 2
)

var orderStatus = [3]string{"Pending", "Cancelled", "Paid"}
var productStatus = [4]string{"Instock", "Low Stock", "Out of Stock", "Preorder"}

type ProductForm struct {
	name        textinput.Model
	sku         textinput.Model
	price       textinput.Model
	categoryId  textinput.Model
	description textarea.Model

	focus int
}

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

	buttonStyle = lipgloss.NewStyle().
			Inherit(baseStyle).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 1).
			BorderForeground(lipgloss.Color("#FF5F87"))

	buttonBlurredStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				Bold(false).
				Padding(0, 1).
				BorderForeground(lipgloss.Color("#353533"))

	divStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(0, 2).
			BorderForeground(lipgloss.Color("69"))
	divBlurredStyle = divStyle.BorderForeground(lipgloss.Color("8"))

	divItem        = baseStyle.Foreground(lipgloss.Color("#FFD046")).Render
	divBlurredItem = baseStyle.Foreground(lipgloss.Color("8")).Render

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
	form       ProductForm
	focusIndex int
	product    models.Product
	response   string
	timer      timer.Model
	cursorMode cursor.Mode
}

func initialModel(product *models.Product) model {
	m := model{
		focusIndex: 0,
		cursorMode: cursor.CursorBlink,
		product:    *product,
		response:   "",
		timer:      timer.New(0),
	}

	var formModel ProductForm
	formModel.focus = 0

	nameInput := textinput.New()
	nameInput.Width = FORM_WIDTH
	nameInput.CharLimit = 70
	nameInput.CursorStyle = cursorStyle
	nameInput.SetValue(product.Name)
	nameInput.Focus()
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle
	formModel.name = nameInput

	skuInput := textinput.New()
	skuInput.Width = FORM_WIDTH
	skuInput.CharLimit = 70
	skuInput.CursorStyle = cursorStyle
	skuInput.SetValue(product.Sku)
	formModel.sku = skuInput

	ta := textarea.New()
	ta.SetHeight(3)
	ta.SetWidth(60)
	ta.SetValue(product.Description)
	ta.SetValue(product.Description)
	formModel.description = ta

	priceInput := textinput.New()
	priceInput.Width = FORM_WIDTH
	priceInput.CharLimit = 70
	priceInput.CursorStyle = cursorStyle
	priceInput.SetValue(fmt.Sprintf("%.2f", product.Price/100))
	formModel.price = priceInput

	catInput := textinput.New()
	catInput.Width = FORM_WIDTH
	catInput.CharLimit = 70
	catInput.CursorStyle = cursorStyle
	catInput.SetValue(fmt.Sprintf("%d", product.CategoryId))
	formModel.categoryId = catInput

	m.form = formModel

	return m
}

func (f *ProductForm) focused() any {
	switch f.focus {
	case 0:
		return &f.name
	case 1:
		return &f.sku
	case 2:
		return &f.description
	case 3:
		return &f.price
	case 4:
		return &f.categoryId
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.timer.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.timer, cmd = m.timer.Update(msg)

	switch msg := msg.(type) {
	case timer.TimeoutMsg:
		m.response = ""
		return m, nil

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.timer.Timeout = timeout
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == MAX_INDEX {
				//Save handler
				price, _ := strconv.ParseFloat(m.form.price.Value(), 64)
				price_in_cents := int(price * 100)
				categoryId, _ := strconv.Atoi(m.form.categoryId.Value())
				_, err := models.UpdateProduct(m.product.Id,
					m.form.name.Value(),
					m.form.sku.Value(),
					m.form.description.Value(),
					price_in_cents,
					categoryId)

				if err != nil {
					m.response = fmt.Sprintf("Error encountered %s", err)
				} else {
					m.response = "Update Successful!"
				}

				return m, m.timer.Start()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > MAX_INDEX {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = MAX_INDEX
			}

			// Remove focused state
			m.form.name.Blur()
			m.form.name.PromptStyle = noStyle
			m.form.name.TextStyle = noStyle

			m.form.sku.Blur()
			m.form.sku.PromptStyle = noStyle
			m.form.sku.TextStyle = noStyle

			m.form.description.Blur()
			//m.form.description.PromptStyle = noStyle
			//m.form.description.TextStyle = noStyle

			m.form.price.Blur()
			m.form.price.PromptStyle = noStyle
			m.form.price.TextStyle = noStyle

			m.form.categoryId.Blur()
			m.form.categoryId.PromptStyle = noStyle
			m.form.categoryId.TextStyle = noStyle

			m.form.focus = m.focusIndex
			focusedInput := m.form.focused()
			var keypressCmd tea.Cmd

			switch fi := focusedInput.(type) {
			case *textinput.Model:
				keypressCmd = fi.Focus()
				fi.PromptStyle = focusedStyle
				fi.TextStyle = focusedStyle
			case *textarea.Model:
				keypressCmd = fi.Focus()
			}

			return m, keypressCmd
		}

	}
	// Handle character input and blinking
	focusedInput := m.form.focused()
	var focusedCmd tea.Cmd

	switch fi := focusedInput.(type) {
	case *textinput.Model:
		var updatedInput textinput.Model
		updatedInput, focusedCmd = fi.Update(msg)
		// Now, update the correct field in the form
		switch m.focusIndex {
		case 0:
			m.form.name = updatedInput
		case 1:
			m.form.sku = updatedInput
		case 3:
			m.form.price = updatedInput
		case 4:
			m.form.categoryId = updatedInput
		}
	case *textarea.Model:
		var updatedInput textarea.Model
		updatedInput, focusedCmd = fi.Update(msg)
		m.form.description = updatedInput
	}

	return m, tea.Batch(cmd, focusedCmd)
}

func formView(form *ProductForm) string {
	focusedEl := form.focused()
	switch fv := focusedEl.(type) {
	case *textinput.Model:
		return fv.View()
	case *textarea.Model:
		return fv.View()
	default:
		return "blank"
	}
}

func (m model) View() string {
	doc := strings.Builder{}
	//physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	{
		w := lipgloss.Width
		leftStatus := statusStyle.Render("<<<<")
		rightStatus := statusStyle.Render(">>>>")
		statusVal := statusText.
			Width(width - w(leftStatus) - w(rightStatus) - 1).Render(fmt.Sprintf("Product ID: %d --- %d %d %d", m.product.Id, m.focusIndex, m.form.focus, int(m.product.Status)))

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			leftStatus,
			statusVal,
			rightStatus,
		)

		serverMessage := baseStyle.Align(lipgloss.Center).Render(m.response)

		div := divStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				divItem(fmt.Sprintf("%s%s", "Name:", m.form.name.View())),
				divItem(fmt.Sprintf("%s%s", "SKU:", m.form.sku.View())),
				divItem(fmt.Sprintf("%s%s", "Description:", m.form.description.View())),
				divItem(fmt.Sprintf("%s%s", "Price:", m.form.price.View())),
				divItem(fmt.Sprintf("%s%s", "CategoryId:", m.form.categoryId.View())),
			),
		)

		var ps strings.Builder
		var status string
		for i := range productStatus {
			if i == int(m.product.Status) {
				status = "x"
				ps.WriteString(divItem(fmt.Sprintf("[%s] %s     ", status, productStatus[i])))
			} else {
				status = " "
				ps.WriteString(divBlurredItem(fmt.Sprintf("[%s] %s     ", status, productStatus[i])))
			}
		}
		statusDiv := divBlurredStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Top,
				ps.String(),
			),
		)

		button := buttonBlurredStyle.Render("Save")
		if m.focusIndex == MAX_INDEX {
			button = buttonStyle.Render("Save")
		}
		//fmt.Fprintf(&b, "\n\n%s\n\n", button)
		body := lipgloss.JoinVertical(lipgloss.Left, bar, serverMessage, div, statusDiv, button)

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
