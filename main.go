package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Categories"})
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

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("categories", getCategories)
		v1.GET("products", getProducts)
		v1.GET("categories/:id", getCategoryById)
		v1.GET("products/:id", getProductById)
	}

	r.Run()
}
