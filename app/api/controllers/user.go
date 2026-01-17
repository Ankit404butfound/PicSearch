package controllers

import (
	"PicSearch/app/api/schemas"
	"PicSearch/app/api/services"
	"fmt"
	"net/http"
	"strconv"

	"PicSearch/app/db/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// GetAllUser
func GetAllUser(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} schemas.UserResponse
// @Failure 404 {object} map[string]string "error": "User not found"
// @Failure 500 {object} map[string]string "error": "Internal server error"
// @Router /users/{id} [get]
// @Param Authorization header string true "Bearer token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE4MDAwNzY0MTUsInVzZXJfaWQiOjF9.CLLSMQGyjT59PRZh1Vx9kdt0uGAcQEisEkFPQkZJzJ4)
func GetUser(c *gin.Context) {
	id := c.Param("id")

	// Convert ID from string to int
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get user details from service
	user, err := services.GetUserByID(userId)
	fmt.Println("user service response:", user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userResponse schemas.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body schemas.CreateUserRequest true "User data"
// @Success 201 {object} schemas.UserResponse
// @Failure 400 {object} map[string]string "error": "Invalid data"
// @Failure 500 {object} map[string]string "error": "Could not create user"
// @Router /users/ [post]
// @Param Authorization header string true "Bearer token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE4MDAwNzY0MTUsInVzZXJfaWQiOjF9.CLLSMQGyjT59PRZh1Vx9k
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Call the service to create the user
	createdUser, err := services.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	var userResponse schemas.UserResponse
	err = copier.Copy(&userResponse, &createdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

// Update user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Convert ID from string to int
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedData models.User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Call the service to update the user
	updatedUser, err := services.UpdateUser(userId, updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE requests to delete a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Convert ID from string to int
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Call the service to delete the user
	err = services.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
