package apiJson

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	// List API endpoints
	router.GET("/api/jewelry", getJewelry)
	router.POST("/api/jewelry", createJewelry)
	router.GET("/api/jewelry/:id", getJewelryById)
	router.PUT("/api/jewelry/:id", updateJewelry)
	router.DELETE("/api/jewelry/:id", deleteJewelry)

	// Jalankan server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
