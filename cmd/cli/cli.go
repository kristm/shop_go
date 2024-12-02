package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"shop_go/internal/config"
	"shop_go/internal/models"
	"strconv"
	"strings"
)

func readCSVFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseCSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	return reader, nil
}

func trim(sliced []string) string {
	output := strings.Join(sliced, "")
	return output
}

func processCSV(reader *csv.Reader) {
	for i := 0; true; i++ {
		if i == 0 {
			continue
		}
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			break
		}

		categoryId, err := strconv.Atoi(trim(record[3:4]))
		if err != nil {
			log.Println("error ", err)
		}
		price, err := strconv.ParseFloat(trim(record[4:5]), 64)
		if err != nil {
			log.Println("error ", err)
		}

		product := models.Product{
			Id:          0,
			Sku:         trim(record[0:1]),
			Name:        trim(record[1:2]),
			Description: trim(record[2:3]),
			CategoryId:  categoryId,
			Price:       price,
		}
		qtyVal := strings.Join(record[5:6], "")
		qty, _ := strconv.Atoi(qtyVal)

		if !readonly {
			productId, err := models.AddProductWithQty(product, qty)
			if err != nil {
				fmt.Printf("ERROR Adding Product: %d %s\n", productId, err)
			}
		}

		fmt.Printf("%d product: %+v %d\n", i, product, qty)
	}
}

func showProducts() {
	products, err := models.GetAllProducts()
	log.Printf("products %v", products)

	if err != nil {
		log.Printf("error getting products %v", err)
	}
	for _, product := range products {
		fmt.Printf("%v\n", product)
	}
}

var readonly = false

func main() {
	showcsv := flag.Bool("showcsv", false, "a bool")
	getproducts := flag.Bool("getproducts", false, "a bool")
	csvPath := flag.String("csv", "", "path to csv")
	flag.Parse()
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config ", err)
	}
	err = models.ConnectDatabase(&cfg)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("ARGC %d %t", len(os.Args), *showcsv)

	if *showcsv {
		readonly = true
	}

	if *getproducts {
		fmt.Println("SHOW PRODUCTS")
		showProducts()
		os.Exit(1)
	}

	data, err := readCSVFile(*csvPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	reader, err := parseCSV(data)
	if err != nil {
		fmt.Println("Error creating CSV reader:", err)
		return
	}
	processCSV(reader)
}
