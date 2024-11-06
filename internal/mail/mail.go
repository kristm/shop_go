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

func NotifyOrder(order *models.Order, customer *models.Customer, cfg *config.Config) (bool, error) {

	t := template.New("template.html")

	var err error
	t, err = t.ParseFiles("template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, order); err != nil {
		log.Println(err)
	}

	emailBody := tpl.String()

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EMAIL_FROM)
	m.SetHeader("To", cfg.EMAIL_REPORTS)
	//m.SetAddressHeader("Cc", "alt email")
	m.SetHeader("Subject", fmt.Sprintf("New Order: %s %s %s", customer.FirstName, customer.LastName, customer.Email))
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.EMAIL_FROM, cfg.EMAIL_PASSWORD)

	if err = d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}
