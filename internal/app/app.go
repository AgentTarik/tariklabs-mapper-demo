package app

import (
	"github.com/gin-gonic/gin"

	"go-test/go-test/internal/core/handler"
	"go-test/go-test/internal/core/repository"
	"go-test/go-test/internal/core/service"
	"go-test/go-test/internal/engine"
)

type App struct {
	router *gin.Engine
}

func New() *App {
	router := gin.Default()

	httpEngine := engine.NewHTTPEngine()
	userRepo := repository.NewUserRepository(httpEngine)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router.GET("/users/:id", userHandler.GetUserByID)
	router.POST("/users", userHandler.CreateUser)

	return &App{
		router: router,
	}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
