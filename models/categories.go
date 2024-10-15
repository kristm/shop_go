package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "/Users/krist/code/shop_go/shop_test.db")
	if err != nil {
		return err
	}
	DB = db
	log.Printf("Connecting to DB")
	log.Printf("DB %+v", db)
	return nil
}

func GetCategories() ([]Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	categories := make([]Category, 0)

	for rows.Next() {
		category := Category{}
		err = rows.Scan(&category.Id, &category.Name)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return categories, err
}

func main() {
	fmt.Println("vim-go")
}
