package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"ticketApp/src/config"
	"ticketApp/src/handler"
	"ticketApp/src/repository"
	"ticketApp/src/service"
)

func main() {
	mCfg := config.NewMongoConfig()
	client, ctx, cancel := mCfg.ConnectDatabase()
	collection := mCfg.GetCollection(client, "User")

	fmt.Println(ctx)
	defer cancel()

	e := echo.New()
	userRepository := repository.NewUserRepository(collection)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	userGroup := e.Group("/api/users")
	userGroup.GET("/:id", userHandler.UserGetById)
	userGroup.POST("", userHandler.UserUpsert)
	userGroup.DELETE("/:id", userHandler.UserDeleteById)

	log.Fatal(e.Start(":8083"))
}
