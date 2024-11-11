package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
)

var DB *sql.DB

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	return func(tb testing.TB) {
		log.Println("teardown suite")
		testTables := []string{"customers", "orders", "order_products", "shipping"}
		log.Println("Models Teardown")
		for _, table := range testTables {
			_, err := ClearTestTable(table)
			if err != nil {
				log.Printf("Teardown error %v", err)
			}
		}
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func ConnectTestDatabase() {
	cfg, err := config.LoadConfig("../.env")
	if err != nil {
		log.Fatal("cannot load config ", err)
	}
	db, err := sql.Open("sqlite3", cfg.TEST_DB)
	if err != nil {
		log.Println(err)
	}
	DB = db
}

func ClearTestTable(tableName string) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	sql := fmt.Sprintf("DELETE FROM %s WHERE id >= ?", tableName)
	stmt, err := DB.Prepare(sql)
	if err != nil {
		return false, err
	}

	defer stmt.Close()
	_, err = stmt.Exec(1)
	if err != nil {
		return false, err
	}
	tx.Commit()
	log.Printf("%s table cleared", tableName)
	return true, nil
}

func TestMain(m *testing.M) {
	log.Println("Test Main")
	ConnectTestDatabase()
	code := m.Run()
	//_, err := ClearProducts()
	//if err != nil {
	//	log.Printf("Teardown error %v\n", err)
	//}
	os.Exit(code)
}

func TestPing(t *testing.T) {
	mailerMock := func(*models.Order, *models.Customer, *config.Config) (bool, error) {
		return true, nil
	}
	router := setupRouter(mailerMock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPostOrders(t *testing.T) {
	mailerMock := func(*models.Order, *models.Customer, *config.Config) (bool, error) {
		return true, nil
	}
	router := setupRouter(mailerMock)

	w := httptest.NewRecorder()
	orders := make([]models.OrderItem, 0)
	order := models.OrderItem{
		ProductId: 1,
		Qty:       2,
		Price:     250.00,
	}

	orders = append(orders, order)

	payload := OrderPayload{
		Orders: orders,
		Customer: models.Customer{
			FirstName: "joe",
			LastName:  "book",
			Email:     "joe@bo.ok",
			Phone:     "123-123",
		},
		Shipping: models.Shipping{
			Address: "malugay st",
		},
	}

	orderJson, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/orders", strings.NewReader(string(orderJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPostIncompleteOrders(t *testing.T) {
	mailerMock := func(*models.Order, *models.Customer, *config.Config) (bool, error) {
		return true, nil
	}
	router := setupRouter(mailerMock)

	w := httptest.NewRecorder()
	orders := make([]models.OrderItem, 0)
	order := models.OrderItem{
		ProductId: 1,
		Qty:       2,
		Price:     250.00,
	}

	orders = append(orders, order)

	payload := OrderPayload{
		Orders:   orders,
		Customer: models.Customer{},
		Shipping: models.Shipping{
			Address: "malugay st",
		},
	}

	orderJson, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/orders", strings.NewReader(string(orderJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
