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
	"golang.org/x/term"
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
	t  table.Model
	vp viewport.Model

	normal = lipgloss.Color("#EEEEEE")
	base   = lipgloss.NewStyle().Foreground(normal)
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#585858"}

	background = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#CCCCCC"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#61D4C6"}

	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			Align(lipgloss.Center)
	dialogTitleStyle = lipgloss.NewStyle().
				Inherit(titleBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type modalModel struct {
	visible bool
	content string
}

type model struct {
	cursor     int
	sections   []string
	rowIndex   int
	tableModel table.Model
	tableRows  int
	targetCol  string
	modal      modalModel
}

func initialModel() model {
	defaultModal := modalModel{
		visible: false,
		content: "hello",
	}
	ordersTable, numRows, columnName := GetOrders(0)
	return model{
		sections:   []string{"Orders", "Customers", "Addresses", "Products", "Vouchers"},
		rowIndex:   0,
		tableModel: ordersTable,
		tableRows:  numRows,
		targetCol:  columnName,
		modal:      defaultModal,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func BlankTable() (table.Model, int, string) {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 10),
	}
	rows := []table.Row{}
	t := table.New(columns).
		WithRows(rows).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true))

	return t, 0, columnKeyID
}

func (m *model) updateTableModel() {
	switch m.cursor {
	case 0:
		m.tableModel, m.tableRows, m.targetCol = GetOrders(m.rowIndex)
	case 1:
		m.tableModel, m.tableRows, m.targetCol = GetCustomers(m.rowIndex)
	case 2:
		m.tableModel, m.tableRows, m.targetCol = GetAddresses(m.rowIndex)
	case 3:
		m.tableModel, m.tableRows, m.targetCol = GetProducts(m.rowIndex)
	case 4:
		m.tableModel, m.tableRows, m.targetCol = GetVouchers(m.rowIndex)
	default:
		m.tableModel, m.tableRows, m.targetCol = BlankTable()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.updateTableModel()
	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmds = append(cmds, tea.Quit)
		case "left", "<":
			if !m.modal.visible {
				MoveLeft(&m)
			}
		case "right", ">":
			if !m.modal.visible {
				MoveRight(&m)
			}
		case "up", "k":
			if !m.modal.visible {
				MoveUp(&m)
			}
		case "down", "j":
			if !m.modal.visible {
				MoveDown(&m)
			}
		case "enter":
			ToggleModal(&m)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

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

		if !m.modal.visible {
			// Tabs
			// get model.rowIndex to determine activeTab
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

			// Table
			{
				m.updateTableModel()
				vp := viewport.New(width, 20)
				vp.YPosition = 20
				vp.SetContent(m.tableModel.View())
				doc.WriteString(vp.View())
			}
		} else {
			titleStr := fmt.Sprintf("INFO %d", m.rowIndex)
			dialogTitle := dialogTitleStyle.Width(physicalWidth - innerMargin).Render(titleStr)
			titleUi := lipgloss.JoinHorizontal(lipgloss.Center, dialogTitle)

			data := GetOrderDetail()
			orderDetail := strings.Join((*data)[0], ": ")
			detailLine := dialogTitleStyle.Render(orderDetail)
			detail := lipgloss.JoinHorizontal(lipgloss.Center, detailLine)

			body := lipgloss.JoinVertical(lipgloss.Center, m.tableModel.HighlightedRow().Data[m.targetCol].(string), detail)

			ui := lipgloss.JoinVertical(lipgloss.Center, body)
			view := lipgloss.JoinVertical(lipgloss.Center, titleUi, ui)
			dialogBody := dialogBoxStyle.Render(view)
			dialog := lipgloss.Place(width, 15,
				lipgloss.Center, lipgloss.Center,
				dialogBody,
				lipgloss.WithWhitespaceChars("商店"),
				lipgloss.WithWhitespaceForeground(subtle),
			)
			doc.WriteString(dialog + "\n\n")
		}

	}

	return docStyle.Render(doc.String())
}

func GetOrders(rowIndex int) (table.Model, int, string) {
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
		table.NewColumn(columnKeyAmount, "Amount", 10),
		table.NewColumn(columnKeyStatus, "Status", 10),
		table.NewColumn(columnKeyCreatedAt, "Created", 25),
	}
	rows := []table.Row{}

	for _, order := range orders {
		amount := fmt.Sprintf("%.2f", order.Amount/100.00)
		status := orderStatus[int(order.Status)]
		customerId := strconv.Itoa(order.CustomerId)
		newRow := table.NewRow(table.RowData{
			columnKeyReference: order.ReferenceCode,
			columnKeyCustomer:  customerId,
			columnKeyAmount:    amount,
			columnKeyStatus:    status,
			columnKeyCreatedAt: order.CreatedAt,
		})
		rows = append(rows, newRow)
	}

	t = resetTable(columns, rows, rowIndex)

	return t, len(rows), columnKeyReference
}

func GetCustomers(rowIndex int) (table.Model, int, string) {
	customers, err := models.GetCustomers()
	if err != nil {
		log.Printf("CUSTOMERS ERROR %v", err)
	}

	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 5).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#88f")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyFirstName, "First Name", 15),
		table.NewColumn(columnKeyLastName, "Last Name", 15),
		table.NewColumn(columnKeyEmail, "Email", 20),
		table.NewColumn(columnKeyPhone, "Phone", 15),
	}
	rows := []table.Row{}

	for _, customer := range customers {
		newRow := table.NewRow(table.RowData{
			columnKeyID:        customer.Id,
			columnKeyFirstName: customer.FirstName,
			columnKeyLastName:  customer.LastName,
			columnKeyEmail:     customer.Email,
			columnKeyPhone:     customer.Phone,
		})
		rows = append(rows, newRow)
	}

	t = resetTable(columns, rows, rowIndex)

	return t, len(rows), columnKeyLastName
}

func GetAddresses(rowIndex int) (table.Model, int, string) {
	addresses, err := models.GetShippingAddresses()
	if err != nil {
		log.Printf("ADDRESSES ERROR %v", err)
	}

	columns := []table.Column{
		table.NewColumn(columnKeyCustomer, "Customer Id", 5),
		table.NewColumn(columnKeyAddress, "Address", 20),
		table.NewColumn(columnKeyCity, "City", 10),
		table.NewColumn(columnKeyCountry, "Country", 10),
		table.NewColumn(columnKeyZip, "Zip", 10),
		table.NewColumn(columnKeyPhone, "Phone", 15),
		table.NewColumn(columnKeyNotes, "Notes", 15),
	}
	rows := []table.Row{}

	for _, address := range addresses {
		newRow := table.NewRow(table.RowData{
			columnKeyCustomer: address.CustomerId,
			columnKeyAddress:  address.Address,
			columnKeyCity:     address.City,
			columnKeyCountry:  address.Country,
			columnKeyZip:      address.Zip,
			columnKeyPhone:    address.Phone,
			columnKeyNotes:    address.Notes,
		})
		rows = append(rows, newRow)
	}

	t = resetTable(columns, rows, rowIndex)

	return t, len(rows), columnKeyCity
}

func GetProducts(rowIndex int) (table.Model, int, string) {
	products, err := models.GetAllProducts()
	if err != nil {
		log.Printf("PRODUCTS ERROR %v", err)
	}

	columns := []table.Column{
		table.NewColumn(columnKeyCategory, "Category", 15).WithStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#88f")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyID, "ID", 10),
		table.NewColumn(columnKeyName, "Name", 25),
		table.NewColumn(columnKeySku, "Sku", 15),
		table.NewColumn(columnKeyPrice, "Price", 15),
		table.NewColumn(columnKeyProductStatus, "Inventory", 20),
	}
	rows := []table.Row{}

	for _, product := range products {
		price := fmt.Sprintf("%.2f", product.Price/100.00)
		newRow := table.NewRow(table.RowData{
			columnKeyID:            product.Id,
			columnKeyCategory:      product.Category,
			columnKeyName:          product.Name,
			columnKeySku:           product.Sku,
			columnKeyPrice:         price,
			columnKeyProductStatus: productStatus[int(product.Status)],
		})
		rows = append(rows, newRow)
	}

	t = resetTable(columns, rows, rowIndex)

	return t, len(rows), columnKeySku
}

func GetVouchers(rowIndex int) (table.Model, int, string) {
	vouchers, err := models.GetVouchers()
	if err != nil {
		log.Printf("VOUCHERS ERROR %v", err)
	}

	columns := []table.Column{
		table.NewColumn(columnKeyVoucherType, "Voucher Type", 10),
		table.NewColumn(columnKeyVoucherCode, "Voucher Code", 15),
		table.NewColumn(columnKeyValid, "Valid", 15),
		table.NewColumn(columnKeyRequiredAmount, "Minimum Amount", 15),
		table.NewColumn(columnKeyDiscount, "Discount", 15),
		table.NewColumn(columnKeyExpires, "Expires", 20),
	}
	rows := []table.Row{}

	for _, voucher := range vouchers {
		newRow := table.NewRow(table.RowData{
			columnKeyVoucherType:    voucher.Type,
			columnKeyVoucherCode:    voucher.Code,
			columnKeyValid:          voucher.Valid,
			columnKeyRequiredAmount: voucher.RequiredAmount,
			columnKeyDiscount:       voucher.Amount,
			columnKeyExpires:        voucher.Expires,
		})
		rows = append(rows, newRow)
	}

	t = resetTable(columns, rows, rowIndex)

	return t, len(rows), columnKeyVoucherCode
}

func resetTable(columns []table.Column, rows []table.Row, rowIndex int) table.Model {
	//footer := fmt.Sprintf("rows: %d", len(rows))
	return table.New(columns).
		WithRows(rows).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		Focused(true).
		Border(customBorder).
		WithPageSize(10).
		WithHighlightedRow(rowIndex).
		WithKeyMap(table.KeyMap{}).
		WithBaseStyle(
			lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#a38")).
				Foreground(lipgloss.Color("#a7a")).
				Align(lipgloss.Left),
		)
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
