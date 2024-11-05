package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type OrderStatus int

const (
	Pending OrderStatus = iota
	Cancelled
	Paid
)

type Order struct {
	Id               int         `json:"id"`
	ShippingId       int         `json:"shipping_id"`
	CustomerId       int         `json:"customer_id"`
	ReferenceCode    string      `json:"reference_code"`
	PaymentReference string      `json:"payment_reference"`
	Amount           float64     `json:"amount_in_cents"`
	Status           OrderStatus `json:"status"`
	Items            []OrderItem `json:"orders"`
}

type OrderItem struct {
	Id        int     `json:"id"`
	OrderId   int     `json:"order_id"`
	ProductId int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
}

// TODO
func (order Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	computedAmount := 0.0
	for _, item := range order.Items {
		total := float64(item.Qty) * float64(item.Price)
		computedAmount += total
	}
	return json.Marshal(&struct {
		*Alias
		Amount float64 `json:"amount"`
	}{
		Alias:  (*Alias)(&order),
		Amount: computedAmount,
	})
}

func (order *Order) UnmarshalJSON(p []byte) error {
	type Alias Order
	aux := &struct {
		Amount float64 `json:"amount_in_cents"`
		*Alias
	}{
		Alias: (*Alias)(order),
	}

	if err := json.Unmarshal(p, &aux); err != nil {
		return err
	}

	//calculation of amount is not necessary if we're storing amount in db
	computedAmount := 0.0
	for _, item := range order.Items {
		total := item.Qty * int(item.Price)
		computedAmount += float64(total)
	}

	order.Amount = computedAmount
	return nil
}

func (prod *OrderItem) UnmarshalJSON(p []byte) error {
	type Alias OrderItem
	aux := &struct {
		Price float64 `json:"price"`
		*Alias
	}{
		Alias: (*Alias)(prod),
	}

	if err := json.Unmarshal(p, &aux); err != nil {
		return err
	}

	prod.Price = aux.Price * 100.00
	return nil
}

func generateReferenceCode() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating reference code %v\n", err)
	}

	return fmt.Sprintf("%X", b)
}

func GetOrderByReference(ref string) (Order, error) {
	stmt, err := DB.Prepare("SELECT id, reference_code, payment_reference, amount_in_cents, status FROM orders WHERE reference_code = ?")
	if err != nil {
		return Order{}, err
	}

	order := Order{}
	sqlErr := stmt.QueryRow(ref).Scan(&order.Id, &order.ShippingId, &order.CustomerId, &order.ReferenceCode, &order.PaymentReference, &order.Amount, &order.Status)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Order{}, nil
		}
		return Order{}, sqlErr
	}

	stmt, err = DB.Prepare("SELECT p.name, op.qty, op.price_in_cents FROM order_products as op LEFT JOIN products as p ON p.id = op.product_id WHERE o.order_id = ?")
	if err != nil {
		return Order{}, err
	}

	rows, sqlErr := stmt.Query(order.Id)
	if sqlErr != nil {
		return Order{}, sqlErr
	}

	defer rows.Close()
	orderItems := make([]OrderItem, 0)

	for rows.Next() {
		orderItem := OrderItem{}
		err = rows.Scan(&orderItem.Id, &orderItem.OrderId, &orderItem.ProductId, &orderItem.Qty, &orderItem.Price)

		if err != nil {
			return Order{}, err
		}

		orderItems = append(orderItems, orderItem)
	}
	//handle error
	err = rows.Err()

	order.Items = orderItems

	return order, nil
}

func GetOrders(customerId int) ([]Order, error) {
	stmt, err := DB.Prepare("SELECT id, customer_id, reference_code, amount_in_cents FROM orders WHERE customer_id = ?")
	if err != nil {
		return nil, err
	}

	rows, sqlErr := stmt.Query(customerId)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	orders := make([]Order, 0)

	for rows.Next() {
		order := Order{}
		err = rows.Scan(&order.Id, &order.CustomerId, &order.Amount, &order.Status)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	err = rows.Err()

	return orders, nil
}

func GetOrderItems(orderId int) ([]OrderItem, error) {
	stmt, err := DB.Prepare("SELECT id, order_id, product_id, qty, price_in_cents FROM order_products WHERE order_id = ?")
	if err != nil {
		return nil, err
	}

	rows, sqlErr := stmt.Query(orderId)
	if sqlErr != nil {
		return nil, sqlErr
	}

	defer rows.Close()
	orderItems := make([]OrderItem, 0)

	for rows.Next() {
		item := OrderItem{}
		err = rows.Scan(&item.Id, &item.OrderId, &item.ProductId, &item.Qty, &item.Price)

		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, item)
	}
	err = rows.Err()

	return orderItems, nil
}

func computeTotalAmount(orderItems []OrderItem) float64 {
	computedAmount := 0.0
	for _, item := range orderItems {
		total := item.Qty * int(item.Price)
		computedAmount += float64(total)
	}

	return computedAmount
}

func AddOrder(order Order) (int, string, error) {
	orderId, refCode, err := AddOrderRecord(order)
	if err != nil {
		return -1, "", err
	}
	for _, item := range order.Items {
		item.OrderId = orderId
		_, err := AddOrderItem(item)
		if err != nil {
			return -1, "", err
		}
	}

	return orderId, refCode, nil
}

func AddOrderRecord(newOrder Order) (int, string, error) {
	tx, err := DB.Begin()
	if err != nil {
		return -1, "", err
	}

	stmt, err := tx.Prepare("INSERT INTO orders (customer_id, shipping_id, reference_code, amount_in_cents, status, payment_reference) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		return -1, "", err
	}

	defer stmt.Close()

	newOrder.Amount = computeTotalAmount(newOrder.Items)
	newOrder.ReferenceCode = generateReferenceCode()

	res, err := stmt.Exec(newOrder.CustomerId, newOrder.ShippingId, newOrder.ReferenceCode, newOrder.Amount, newOrder.Status, newOrder.PaymentReference)

	if err != nil {
		return -1, "", err
	}

	tx.Commit()
	orderId, _ := res.LastInsertId()

	//return int(orderId), newOrder.ReferenceCode, nil
	return int(orderId), newOrder.ReferenceCode, nil
}

func AddOrderItem(newOrderItem OrderItem) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	//TODO: fetch price from product table
	stmt, err := tx.Prepare("INSERT INTO order_products (order_id, product_id, qty, price_in_cents) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newOrderItem.OrderId, newOrderItem.ProductId, newOrderItem.Qty, newOrderItem.Price)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
