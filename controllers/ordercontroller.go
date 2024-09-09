package controllers

import (
	"net/http"

	"github.com/cngJo/golang-api-auth/database"
	"github.com/cngJo/golang-api-auth/models"
	"github.com/gin-gonic/gin"
)

// GetOrders retrieves a list of orders
func GetOrders(context *gin.Context) {
	var orders []models.Order
	if err := database.Instance.Preload("Product").Preload("User").Find(&orders).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// UpdateOrderStatus updates the status of an order (processed or completed)
func CreateOrderStatus(context *gin.Context) {
	// Ambil ID produk dari parameter URL
	var orders models.Order // Cari produk berdasarkan ID

	if err := context.ShouldBindJSON(&orders); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan produk ke database
	if err := database.Instance.Create(&orders).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}
