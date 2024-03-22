package main

import (
	"assignment2/db"
	"assignment2/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createOrder(c *gin.Context) {
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func getOrders(c *gin.Context) {
	var orders []models.Order
	db.DB.Preload("Items").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func getOrder(c *gin.Context) {
	var order models.Order
	if err := db.DB.Preload("Items").Where("order_id = ?", c.Param("orderId")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func updateOrder(c *gin.Context) {
	var order models.Order
	if err := db.DB.Preload("Items").Where("order_id = ?", c.Param("orderId")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func deleteOrder(c *gin.Context) {
	var order models.Order
	if err := db.DB.Preload("Items").Where("order_id = ?", c.Param("orderId")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	if err := db.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Record deleted!"})
}

func setupRoutes() *gin.Engine {
	r := gin.Default()
	r.POST("/orders", createOrder)
	r.GET("/orders", getOrders)
	r.GET("/orders/:orderId", getOrder)
	r.PUT("/orders/:orderId", updateOrder)
	r.DELETE("/orders/:orderId", deleteOrder)
	return r
}

func main() {
	db.Connect()
	r := setupRoutes()
	_ = db.DB.AutoMigrate(&models.Order{})
	err := r.Run(":9000")
	if err != nil {
		return
	}
	fmt.Println("Server is running on port 8080")
}
