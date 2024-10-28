package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"shop_go/models"
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

		//if !readonly {
		//	_, err = models.AddProduct(product)
		//	if err != nil {
		//		fmt.Println("ERROR: ", err)
		//	}
		//}

		qty := strings.Join(record[5:6], "")
		//qty, _ = strconv.Atoi(qty)
		fmt.Printf("%d product: %+v %s\n", i, product, qty)
		//fmt.Printf("%s %s %s %s %s %s\n", record[0:1], record[1:2], record[2:3], record[3:4], record[4:5], record[5:6])
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
	flag.Parse()
	err := models.ConnectDatabase()
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("ARGC %d %t", len(os.Args), *showcsv)

	if *showcsv {
		readonly = true
		fmt.Println("SHOW PRODUCTS")
		showProducts()
		os.Exit(1)
	}

	data, err := readCSVFile("mini_inventory.csv")
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
