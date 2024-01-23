package main

import (
	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/EM/internal/handlers"
)

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
