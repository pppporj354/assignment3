package main

import (
	"assignment3/db"
	"assignment3/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("fasdgsdgsd"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

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
	r.Use(AuthMiddleware())
	r.POST("/orders", AuthMiddleware(), createOrder)
	r.GET("/orders", AuthMiddleware(), getOrders)
	r.GET("/orders/:orderId", AuthMiddleware(), getOrder)
	r.PUT("/orders/:orderId", AuthMiddleware(), updateOrder)
	r.DELETE("/orders/:orderId", AuthMiddleware(), deleteOrder)
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
