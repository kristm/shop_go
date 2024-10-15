package main

import (
	"log"
	"net/http"
	"shop_go/models"

	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
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
	c.JSON(http.StatusOK, gin.H{"message": "Products"})
}
func getCategoryById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Category" + id})
}
func getProductById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Product" + id})
}

func getRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("", getRoot)
		v1.GET("categories", getCategories)
		v1.GET("products", getProducts)
		v1.GET("categories/:id", getCategoryById)
		v1.GET("products/:id", getProductById)
	}

	r.Run()
}
