package tui

import (
	"fmt"
	"shop_go/internal/models"
)

type OrderDetail [][]string

func GetOrderDetail(ref string) *OrderDetail {
	var orderDetail OrderDetail
	order, _ := models.GetOrderByReference(ref)

	for _, orderItem := range order.Items {
		qty := fmt.Sprintf("%d", orderItem.Qty)
		price := fmt.Sprintf("%.2f", orderItem.Price/100.00)
		orderDetail = append(orderDetail, []string{orderItem.ProductName, qty, price})
	}

	return &orderDetail
}

type SubscriberList [][]string
