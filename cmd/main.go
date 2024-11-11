package main

import (
	"log"
	"net/http"
	"shop_go/internal/config"
	"shop_go/internal/mail"
	"shop_go/internal/models"
	"strconv"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type OrderPayload struct {
	Orders           []models.OrderItem
	Customer         models.Customer
	Socials          models.Socials
	Shipping         models.Shipping
	PaymentReference string `json:"reference"`
	DeviceProfile    string `json:"device_profile"`
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

func getOrderByReference(c *gin.Context) {
	reference := c.Param("reference_code")
	order, err := models.GetOrderByReference(reference)
	checkErr(err)

	if order != nil {
		c.JSON(http.StatusOK, gin.H{"data": order})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func sendOrderNotification(order *models.Order, customer *models.Customer, cfg *config.Config) error {
	_, err := mail.NotifyOrder(order, customer, cfg)
	if err != nil {
		log.Printf("EMAIL ERROR %v\n", err)
		return err
	}

	return nil
}

func createOrder(c *gin.Context) {
	var requestBody OrderPayload
	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Error parsing POST payload %v\n", err)
	}

	// create customer record
	customerId, err := models.AddOrGetCustomer(&requestBody.Customer)
	if err != nil {
		log.Printf("Error Adding Customer %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	// create social media record
	requestBody.Socials.CustomerId = customerId
	models.AddCustomerSocials(&requestBody.Socials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	// create shipping record
	requestBody.Shipping.CustomerId = customerId
	shippingId, err := models.AddShipping(&requestBody.Shipping)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	// create analytics
	_, err = models.AddAnalytics(&models.Analytics{
		CustomerId: customerId,
		IpAddress:  c.ClientIP(),
		Device:     requestBody.DeviceProfile,
	})
	if err != nil {
		log.Printf("ERROR creating analytics record %v\n", err)
	}

	// create order record
	log.Printf("Order Items Req %+v", requestBody.Orders)
	order := models.Order{
		ShippingId:       shippingId,
		CustomerId:       customerId,
		Status:           0,
		PaymentReference: requestBody.PaymentReference,
		Items:            requestBody.Orders,
	}
	success, referenceCode, err := models.AddOrder(order)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		order.ReferenceCode = referenceCode
		cfg, _ := loadConfig()
		go sendOrderNotification(&order, &requestBody.Customer, cfg)
		requestBody.Customer.Id = customerId
		log.Println(customerId)
		log.Printf("json payload %v\n", requestBody)
		log.Printf("order: %v\n", requestBody.Orders)
		log.Printf("payment reference: %s\n", requestBody.PaymentReference)
		log.Printf("customer: %v\n", requestBody.Customer)
		log.Printf("socials: %v\n", requestBody.Socials)
		log.Printf("shipping: %v\n", requestBody.Shipping)

		c.JSON(http.StatusOK, gin.H{"message": "ok", "reference_code": referenceCode, "success": success})
	}
}

func loadConfig() (*config.Config, error) {
	var configPath string
	if testing.Testing() {
		configPath = "../.env"
	} else {
		configPath = ".env"
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("cannot load config ", err)
		return nil, err
	}

	return &cfg, err
}

func setupRouter() *gin.Engine {
	cfg, _ := loadConfig()
	err := models.ConnectDatabase(cfg)
	checkErr(err)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("categories", getCategories)
		v1.GET("categories/:id", getCategoryById)
		v1.GET("products/category/:category_id", getProducts)
		v1.GET("products/:id", getProductById)
		v1.GET("products/sku/:sku", getProductBySku)
		v1.GET("orders/:reference_code", getOrderByReference)
		v1.POST("orders", createOrder)
	}
	return router
}

func main() {

	r := setupRouter()

	r.Run()
}
