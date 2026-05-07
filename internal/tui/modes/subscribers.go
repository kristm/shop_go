package modes

import (
	"log"
	"shop_go/internal/models"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

func ShowSubscribers() string {
	doc := strings.Builder{}
	subscribers, err := models.GetSubscribers()
	if err != nil {
		log.Printf("Error getting subscribers")
		return ""
	}
	subsList := []table.Row{}
	subsCount := len(subscribers)
	for i := 0; i < subsCount; i++ {
		row := table.NewRow(table.RowData{
			"EMAIL":      subscribers[i].Email,
			"SUBSCRIBED": subscribers[i].CreatedAt,
		})
		subsList = append(subsList, row)
	}
	columns := []table.Column{
		table.NewColumn("EMAIL", "EMAIL", 40).WithStyle(baseStyle),
		table.NewColumn("SUBSCRIBED", "SUBSCRIBE DATE", 40).WithStyle(center),
	}
	t := table.New(columns).
		WithRows(subsList).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
		WithBaseStyle(
			lipgloss.NewStyle().
				Padding(2).
				BorderForeground(cyan),
		)
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("Subscribers"), t.View()))
	return docStyle.Render(doc.String())
}
