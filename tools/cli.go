package main

import (
	"bytes"
	"encoding/csv"
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

		categoryId, err := strconv.Atoi(trim(record[4:5]))
		if err != nil {
			log.Println("error ", err)
		}
		price, err := strconv.ParseFloat(trim(record[5:6]), 64)
		if err != nil {
			log.Println("error ", err)
		}

		product := models.Product{
			Id:          0,
			Sku:         trim(record[0:1]),
			Name:        trim(record[2:3]),
			Description: trim(record[3:4]),
			CategoryId:  categoryId,
			Price:       price,
		}

		_, err = models.AddProduct(product)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		fmt.Printf("%d product: %+v\n", i, product)
	}
}

func main() {
	err := models.ConnectDatabase()
	if err != nil {
		log.Printf("%s\n", err)
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
