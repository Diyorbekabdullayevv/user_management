package main

import (
	"fmt"
	"log"
	"mongodb-project/db"
	"mongodb-project/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8080"

	// 🔑 CONNECT TO MONGODB
	if err := db.Connect("mongodb://localhost:27017"); err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	server := gin.Default()

	// Serve CSS files
	server.Static("/css", "../frontend/css")

	// Serve JavaScript files
	server.Static("/js", "../frontend/js")

	// Frontend route - serve index.html
	server.GET("/", func(c *gin.Context) {
		c.File("../frontend/html/index.html")
	})

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	server.POST("/users", internal.CreateUser)
	server.POST("/api/users", internal.StoreUserData)
	server.GET("/api/users/filter", internal.GetUsersByAge)
	server.POST("/users/bulk", internal.BulkCreateUsers)
	server.GET("/users/:id", internal.GetUserByID)
	server.GET("/users/all", internal.GetUsers)
	server.PUT("/users/update/:id", internal.UpdateUserById)
	server.PUT("/users/update/bulk", internal.UpdateUsers)
	server.DELETE("/users/delete/:id", internal.DeleteUserByID)
	server.GET("/docktor", internal.GetDocktor)

	fmt.Println("Server running on port:", port)
	server.Run(port)
}
