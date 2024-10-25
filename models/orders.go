package models

import (
	"encoding/json"
)

type Status int

const (
	Pending Status = iota
	Cancelled
	Paid
)

type Order struct {
	Id         int         `json:"id"`
	ShippingId int         `json:"shipping_id"`
	CustomerId int         `json:"customer_id"`
	Amount     float64     `json:"amount_in_cents"`
	Status     Status      `json:"status"`
	Items      []OrderItem `json:"orders"`
}

type OrderItem struct {
	OrderId   int     `json:"id"`
	ProductId int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
}

// TODO
func (order Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	computedAmount := 0.0
	for _, item := range order.Items {
		total := item.Qty * int(item.Price)
		computedAmount += float64(total)
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

	order.Amount = computedAmount / 100.00
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

func GetOrders(customerId int) ([]Order, error) {
	stmt, err := DB.Prepare("SELECT id, customer_id, amount_in_cents FROM orders WHERE customer_id = ?")
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

func AddOrder(newOrder Order) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO orders (customer_id, amount_in_cents, status) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newOrder.CustomerId, newOrder.Amount, newOrder.Status)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func AddOrderItem(newOrderItem OrderItem) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO orders_product (order_id, product_id, qty) VALUES (?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newOrderItem.OrderId, newOrderItem.ProductId, newOrderItem.Qty)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
