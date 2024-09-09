package controllers

import (
	"net/http"

	"github.com/cngJo/golang-api-auth/database"
	"github.com/cngJo/golang-api-auth/models"
	"github.com/gin-gonic/gin"
)

func CreateProduct(context *gin.Context) {
	// Ambil data produk dari request JSON
	var product models.Product
	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan produk ke database
	if err := database.Instance.Create(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	// Hanya mengembalikan pesan sukses tanpa rincian produk
	context.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
	})
}

func UpdateProduct(context *gin.Context) {
	// Ambil ID produk dari parameter URL
	var product models.Product
	id := context.Param("id") // Perbaiki menjadi "id"

	// Cari produk berdasarkan ID
	if err := database.Instance.First(&product, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Ambil data dari request JSON dan update produk
	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan perubahan ke database
	if err := database.Instance.Save(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}

	// Berikan respon berhasil
	context.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
	})
}

func DeleteProduct(context *gin.Context) {
	// Ambil ID produk dari parameter URL
	id := context.Param("id")

	// Cari produk berdasarkan ID
	var product models.Product
	if err := database.Instance.Where("id = ?", id).First(&product).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Hapus produk dari database secara permanen menggunakan Unscoped
	if err := database.Instance.Unscoped().Delete(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}

	// Berikan respon berhasil
	context.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
func GetProductDetail(context *gin.Context) {
	// Ambil ID produk dari parameter URL
	id := context.Param("id")

	// Cari produk berdasarkan ID
	var product models.Product
	if err := database.Instance.Where("id = ?", id).First(&product).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Hapus produk dari database secara permanen menggunakan Unscoped
	if err := database.Instance.Unscoped().Where(&product).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error product"})
		return
	}

	// Berikan respon produk yang ditemukan dengan format yang sesuai
	context.JSON(http.StatusOK, gin.H{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"user_id":     product.UserID,
	})
}
