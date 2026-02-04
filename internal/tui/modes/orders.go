package modes

import (
	"fmt"
	"log"
	"os"
	"shop_go/internal/models"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

var (
	orderStatuses = []string{"Pending", "Cancelled", "Paid"}
	shipStatuses  = []string{"Packing", "In transit", "Shipped", "Delivered"}
	doc           = strings.Builder{}

	cyan = lipgloss.Color("#00DBC6")

	//baseStyle = lipgloss.NewStyle().
	//		Bold(true).
	//		Foreground(lipgloss.Color("#FFF7DB"))

	titleStyle = baseStyle.
			Padding(0, 1).
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			MarginBottom(1).
			Background(cyan)

	center = lipgloss.NewStyle().
		Align(lipgloss.Center)
	div = lipgloss.NewStyle().
		Padding(1, 0).
		Width(columnWidth - 10)

	//divItem = baseStyle.Foreground(lipgloss.Color("#FFD046")).Render
	header = baseStyle.
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan).
		Align(lipgloss.Center).
		Width(columnWidth - 12)

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan).
		PaddingLeft(1).
		Width(columnWidth - 10)

	listHeader = baseStyle.
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(cyan).
			Render

	listItem = baseStyle.PaddingLeft(2).Render
	//docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

var columnWidth int

func ShowOrder(reference string) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if physicalWidth > 50 {
		columnWidth = 60
	} else {
		columnWidth = 50
	}

	order, err := models.GetOrderByReference(reference)
	if err != nil {
		log.Printf("Error getting orders %v", err)
		os.Exit(1)
	}

	customer, err := models.GetCustomerById(order.CustomerId)
	if err != nil {
		log.Printf("error getting orders %v", err)
	}

	shipping, err := models.GetShippingById(order.ShippingId)
	if err != nil {
		log.Printf("error getting orders %v", err)
	}

	name := fmt.Sprintf("%s %s", customer.FirstName, customer.LastName)
	phone := fmt.Sprintf("%s", customer.Phone)
	email := fmt.Sprintf("%s", customer.Email)
	address := fmt.Sprintf("%s", shipping.Address)
	city := fmt.Sprintf("%s", shipping.City)
	country_zip := fmt.Sprintf("%s %s", shipping.Country, shipping.Zip)
	notes := fmt.Sprintf("NOTES: %s", shipping.Notes)
	customerDetails := div.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header.Render(fmt.Sprintf("Order %s", reference)),
			titleStyle.Render("Customer Details"),
			divItem(name),
			divItem(phone),
			divItem(email),
			divItem(address),
			divItem(city),
			divItem(country_zip),
			divItem(notes),
		),
	)

	amount := fmt.Sprintf("Total: %.2f", order.Amount/100.00)
	orderId := fmt.Sprintf("ID: %d", order.Id)
	referenceLine := fmt.Sprintf("Reference: %s", reference)
	orderStatus := fmt.Sprintf("Status: %s", orderStatuses[order.Status])
	paymentRef := fmt.Sprintf("Payment Reference: %s", order.PaymentReference)
	voucher := fmt.Sprintf("Voucher: %s", order.VoucherCode)
	customerId := fmt.Sprintf("Customer ID: %d", order.CustomerId)
	shippingId := fmt.Sprintf("Shipping ID: %d", order.ShippingId)
	shippingStatus := fmt.Sprintf("Shipping Status: %s", shipStatuses[order.ShippingStatus])
	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Order Details"),
				listItem(amount),
				listItem(orderId),
				listItem(referenceLine),
				listItem(orderStatus),
				listItem(paymentRef),
				listItem(voucher),
				listItem(customerId),
				listItem(shippingId),
				listItem(shippingStatus),
			),
		),
	)

	items := []table.Row{}
	orderRows := len(order.Items)
	for i := 0; i < orderRows; i++ {
		row := table.NewRow(table.RowData{
			"NAME":  order.Items[i].ProductName,
			"PRICE": fmt.Sprintf(" %.2f ", order.Items[i].Price/100.00),
			"QTY":   fmt.Sprintf(" %d ", order.Items[i].Qty),
		})
		items = append(items, row)
	}
	var nameWidth int
	if physicalWidth > 50 {
		nameWidth = 40
	} else {
		nameWidth = columnWidth / 3
	}
	columns := []table.Column{
		table.NewColumn("NAME", "NAME", nameWidth).WithStyle(baseStyle),
		table.NewColumn("QTY", "QTY", 10).WithStyle(center),
		table.NewColumn("PRICE", "PRICE", 10).WithStyle(center),
	}
	t := table.New(columns).
		WithRows(items).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		WithBaseStyle(
			lipgloss.NewStyle().
				Padding(2).
				BorderForeground(cyan),
		)

	firstRow := lipgloss.JoinHorizontal(lipgloss.Top, customerDetails)
	var secondRow string
	if columnWidth > 50 {
		secondRow = lipgloss.JoinHorizontal(lipgloss.Top, lists, t.View())
	} else {
		secondRow = lipgloss.JoinVertical(lipgloss.Left, lists, t.View())
	}

	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, firstRow, secondRow))
	return docStyle.Render(doc.String())
}
