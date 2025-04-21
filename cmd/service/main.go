package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nikborovets/backend-trainee-assignment-spring-2025/configs"
)

func main() {
	_ = godotenv.Load()
	cfg := configs.LoadConfig()
	_ = cfg // пока не используем, но пригодится для usecase

	r := gin.Default()

	// TODO: инициализация репозиториев, usecase, DI
	// Пока просто healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
