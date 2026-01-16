package controllers

import (
	"PicSearch/app/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param authRequest body AuthRequest true "Login credentials"
// @Success 200 {object} map[string]string "token": "JWT token"
// @Failure 400 {object} map[string]string "error": "Invalid request"
// @Failure 401 {object} map[string]string "error": "Authentication failed"
// @Router /login [post]
func Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
