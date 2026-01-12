package router

import (
	"PicSearch/app/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	// Group user-related routes
	userGroup := apiGroup.Group("/users")
	{
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/", controllers.GetAllUser)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}

	apiGroup.GET("/", HealthCheck)
}

// HealthCheck godoc
// @Summary Health Check
// @Description Returns the status of the API
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "API is running",
	})
}
