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
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating orders"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}

func UpdateOrderStatus(context *gin.Context) {
	// Ambil ID produk dari parameter URL
	var orders models.Order   // Cari produk berdasarkan ID
	id := context.Param("id") // Perbaiki menjadi "id"

	// Cari produk berdasarkan ID
	if err := database.Instance.First(&orders, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "orders not found"})
		return
	}

	// Ambil data dari request JSON dan update produk
	if err := context.ShouldBindJSON(&orders); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan perubahan ke database
	if err := database.Instance.Save(&orders).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating orders"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}
