package models

import "encoding/json"

type Status int

const (
	Pending Status = iota
	Cancelled
	Paid
)

type Order struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customer_id"`
	Amount     int    `json:"amount_in_cents"`
	Status     Status `json:"status"`
}

type OrderItem struct {
	ProductId  int `json:"product_id"`
	CategoryId int `json:"category_id"`
	Qty        int `json:"qty"`
	Price      int `json:"price"`
}

func (prod *OrderItem) UnmarshalJSON(p []byte) error {
	type Alias OrderItem
	aux := &struct {
		Price int `json:"price"`
		*Alias
	}{
		Alias: (*Alias)(prod),
	}

	if err := json.Unmarshal(p, &aux); err != nil {
		return err
	}

	//TODO: convert Price from float to int
	prod.Price = aux.Price * 100
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
