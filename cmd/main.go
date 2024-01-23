package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"EM/internal/handlers"
)

func init() {
    // loads values from .env into the system
    if err := godotenv.Load("dsn.env"); err != nil {
        log.Print("No .env file found")
    }
}

func main() {

	newApp := handlers.NewApp()

	// задаем роутер и хендлеры
	router := gin.Default()
	router.POST("/user", newApp.AddUser)
	router.DELETE("/user/:id", newApp.RemoveUser)
	router.PATCH("/user/:id", newApp.UpdateUser)
	router.GET("/users", newApp.GetUsers)

	// запускаем сервер
	router.Run("localhost:8080")
}
