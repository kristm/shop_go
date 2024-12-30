package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"shop_go/internal/config"
	"shop_go/internal/models"

	gomail "github.com/Shopify/gomail"
)

func StatusLabel(status models.OrderStatus) string {
	return [...]string{"Pending", "Cancelled", "Paid"}[status]
}

func NotifyOrder(order *models.Order, customer *models.Customer, cfg *config.Config) (bool, error) {

	var err error
	t, err := template.New("template.html").Funcs(template.FuncMap{
		"StatusLabel": StatusLabel,
	}).ParseFiles("internal/mail/template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer

	//convert from cents
	//TODO how to use existing json marshalling
	computedAmount := 0.0
	for i := 0; i < len(order.Items); i++ {
		p := order.Items[i].Price
		order.Items[i].Price = p / 100.00
		computedAmount += order.Items[i].Price * float64(order.Items[i].Qty)
	}

	order.Amount = computedAmount

	if err = t.Execute(&tpl, order); err != nil {
		log.Println(err)
	}

	emailBody := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EMAIL_FROM)
	m.SetHeader("To", customer.Email)
	m.SetAddressHeader("Cc", cfg.EMAIL_REPORTS)
	m.SetHeader("Subject", fmt.Sprintf("New Order: %s %s %s", customer.FirstName, customer.LastName, customer.Email))
	m.SetBody("text/html", emailBody)
	log.Printf("EMAIL %v", emailBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.EMAIL_FROM, cfg.EMAIL_PASSWORD)
	d.StartTLSPolicy = gomail.MandatoryStartTLS

	if err = d.DialAndSend(m); err != nil {
		log.Printf("MAIL ERROR %v", err)
		return false, err
	}

	return true, nil
}
