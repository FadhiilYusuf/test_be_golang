package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/cngJo/golang-api-auth/database"
	"github.com/cngJo/golang-api-auth/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Secret Key (should be stored in an environment variable)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Function to generate a JWT token
func generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(400, gin.H{"error": "Error hashing password"})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(400, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	// Generate JWT Token
	token, err := generateJWT(user)
	if err != nil {
		context.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	// Return user data including encrypted password
	context.JSON(200, gin.H{
		"username": user.Username,
		"token":    token,
		"password": user.Password, // Encrypted password returned in response
	})
}

// Function to handle user login
func LoginUser(context *gin.Context) {
	var credentials struct {
		Username string `gorm:"unique"`
		Password string `json:"password"`
	}

	if err := context.ShouldBindJSON(&credentials); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.Instance.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		// If the username is not found, return a failure message
		context.JSON(400, gin.H{"error": "Invalid username or password", "message": "Login failed"})
		return
	}

	if err := user.CheckPassword(credentials.Password); err != nil {
		// If the password doesn't match, return a failure message
		context.JSON(400, gin.H{"error": "Invalid username or password", "message": "Login failed"})
		return
	}

	// Return a success message with the token
	context.JSON(200, gin.H{
		"message": "Login successful",
	})
}

// Function to reset the user's password
func ResetPassword(context *gin.Context) {
	var request struct {
		Username string `gorm:"unique"`
		Password string `json:"password"`
	}

	// Bind the JSON request to the request struct
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error(), "message": "Invalid request format"})
		return
	}

	var user models.User

	// Find the user by username
	if err := database.Instance.Where("username = ?", request.Username).First(&user).Error; err != nil {
		context.JSON(404, gin.H{"error": "User not found", "message": "Reset password failed"})
		return
	}

	// Hash the new password
	if err := user.HashPassword(request.Password); err != nil {
		context.JSON(500, gin.H{"error": "Error hashing password", "message": "Reset password failed"})
		return
	}

	// Update the password in the database
	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(500, gin.H{"error": "Error saving new password", "message": "Reset password failed"})
		return
	}

	// Return success message with the token
	context.JSON(200, gin.H{
		"message": "Password reset successful",
	})
}

func LogoutUser(context *gin.Context) {
	// Inform the client to handle token removal
	context.JSON(200, gin.H{
		"message": "Logout successful",
	})
}

func UpdateProfile(context *gin.Context) {
	id := context.Param("id")

	var updateData struct {
		Username string `json:"username"`
		Foto     string `json:"foto"`
	}

	if err := context.ShouldBindJSON(&updateData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.Instance.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user fields
	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Foto != "" {
		user.Foto = updateData.Foto
	}

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}
