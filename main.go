package main

import (
	"log"
	"net/http"
	"shop_go/models"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type OrderPayload struct {
	Orders   []models.OrderItem
	Customer models.Customer
	Shipping models.Shipping
}

func checkErr(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}

func getCategories(c *gin.Context) {
	categories, err := models.GetCategories()
	checkErr(err)

	if categories == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": categories})
	}
}

func getProducts(c *gin.Context) {
	category_id, err := strconv.Atoi(c.Param("category_id"))
	checkErr(err)
	products, err := models.GetProducts(category_id)
	checkErr(err)

	if products == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": products})
	}
}

func getCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	category, err := models.GetCategoryById(id)
	checkErr(err)

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": category})
	}
}

func getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	checkErr(err)
	product, err := models.GetProductById(id)
	checkErr(err)

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": product})
	}
}

func getProductBySku(c *gin.Context) {
	sku := c.Param("sku")
	product, err := models.GetProductBySku(sku)
	checkErr(err)

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No records found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": product})
	}
}

func createOrder(c *gin.Context) {
	// create customer record
	// create shipping record
	// create order record
	// create order products join table entries
	var requestBody OrderPayload

	// ERROR: Error parsing POST payload json: cannot unmarshal number 349.99 into Go struct field OrderItem.Orders.price of type int
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Error parsing POST payload %v", err)
	}

	log.Printf("json payload %v\n  %v", requestBody, requestBody.Orders)

	c.JSON(http.StatusOK, gin.H{"message": "TODO"})
}

func getRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/api/categories", getCategories)

	v1 := r.Group("/api/v1")
	{
		v1.GET("", getRoot)
		v1.GET("categories", getCategories)
		v1.GET("categories/:id", getCategoryById)
		v1.GET("products/category/:category_id", getProducts)
		v1.GET("products/:id", getProductById)
		v1.GET("products/sku/:sku", getProductBySku)
		v1.POST("orders", createOrder)
	}

	r.Run()
}
