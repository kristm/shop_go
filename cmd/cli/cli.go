package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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

func PrintJSON(obj interface{}) {
	bytes, _ := json.MarshalIndent(obj, "\t", "\t")
	fmt.Println(string(bytes))
}

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

func getOrderDetails(reference string) {
	order, err := models.GetOrderByReference(reference)
	if err != nil {
		log.Printf("error getting orders %v", err)
	}

	customer, err := models.GetCustomerById(order.CustomerId)
	if err != nil {
		log.Printf("error getting orders %v", err)
	}

	shipping, err := models.GetShippingById(order.ShippingId)
	if err != nil {
		log.Printf("error getting orders %v", err)
	}

	fmt.Printf("%s %s\n", customer.FirstName, customer.LastName)
	fmt.Printf("%s \n", customer.Phone)
	fmt.Printf("%s \n", customer.Email)
	fmt.Printf("%s \n", shipping.Address)
	fmt.Printf("%s \n", shipping.City)
	fmt.Printf("%s %s\n", shipping.Country, shipping.Zip)
	fmt.Printf("NOTES: %s \n", shipping.Notes)
	PrintJSON(order)
	PrintJSON(shipping)
}

func updateShipping(reference string) {
	order, err := models.GetOrderByReference(reference)
	if err != nil {
		log.Printf("error getting order %v", err)
	}

	status := order.ShippingStatus
	shippingId := order.ShippingId

	if status < models.Delivered {
		status++
		_, err := models.UpdateShippingStatus(shippingId, status)
		if err != nil {
			log.Printf("error updating shipping status %v", err)
		}

		fmt.Printf("Updated Shipping Status: %d\n", status)
	}
}

func markPaidOrder(reference string) {
	//Get Order By reference
	order, _ := models.GetOrderByReference(reference)
	// update Order status
	_, err := models.UpdateOrderStatus(reference, models.Paid)
	if err != nil {
		log.Printf("ERROR: Updating Order Status, Order %s", reference)
		return
	}
	// get Order Items
	// update product inventory for each item
	items, _ := models.GetOrderItems(order.Id)
	for _, item := range items {
		// get item qty
		qty := item.Qty
		productId := item.ProductId
		// get product qty
		inventory, _ := models.GetProductInventory(productId)
		prodQty := inventory.Qty
		// subtract item qty from product qty
		updatedQty := prodQty - qty
		// update product inventory
		_, err := models.UpdateProductInventory(productId, updatedQty)
		if err != nil {
			log.Printf("ERROR: Updating Inventory for Product %d %d", productId, updatedQty)
		}
	}
	fmt.Printf("Set Order as Paid\n")
}

func setPreorderFromSku(sku string) {
	product, err := models.GetProductBySku(sku)
	if err != nil {
		log.Printf("Product does not exist %s", sku)
		return
	}

	_, _ = models.UpdateProductInventory(product.Id, 0)

	ok, err := models.SetPreorder(&product)
	if err != nil {
		log.Printf("Error updating status %s", err)
	}
	if ok {
		log.Printf("Updated Product Status to Preorder: %d", product.Id)
		return
	}
}

var readonly = false

func main() {
	showcsv := flag.Bool("showcsv", false, "a bool")
	getproducts := flag.Bool("getproducts", false, "a bool")
	csvPath := flag.String("csv", "", "path to csv")
	orderRef := flag.String("setpaid", "", "order reference")
	orderRefCode := flag.String("getorder", "", "order reference")
	shippingRef := flag.String("updateship", "", "order reference")
	preorder := flag.String("setpreorder", "", "product sku")
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

	if len(*orderRef) > 0 {
		markPaidOrder(*orderRef)
		os.Exit(1)
	}

	if len(*orderRefCode) > 0 {
		getOrderDetails(*orderRefCode)
		os.Exit(1)
	}

	if len(*shippingRef) > 0 {
		updateShipping(*shippingRef)
		os.Exit(1)
	}

	if len(*preorder) > 0 {
		setPreorderFromSku(*preorder)
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
