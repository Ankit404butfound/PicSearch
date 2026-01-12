package main

import (
	"PicSearch/app/api/router"

	"PicSearch/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // This is the gin-swagger package
)

func main() {
	// Setup Swagger
	docs.SwaggerInfo.Title = "PicSearch API"
	docs.SwaggerInfo.Description = "API documentation for PicSearch"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Database connection string
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run() // Start the server
}
