package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"ticketApp/src/config"
	"ticketApp/src/handler"
	"ticketApp/src/repository"
	"ticketApp/src/service"
	"ticketApp/src/type/util"
)

func init() {

}

func main() {
	mCfg := config.NewMongoConfig()
	client, ctx, cancel := mCfg.ConnectDatabase()
	collection := mCfg.GetCollection(client, "User")

	fmt.Println(ctx)
	defer cancel()

	e := echo.New()
	e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

	userRepository := repository.NewUserRepository(collection)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	userGroup := e.Group("/api/users")
	userGroup.GET("", userHandler.UserGetById)
	userGroup.GET("/:id", userHandler.UserGetById)
	userGroup.POST("", userHandler.UserUpsert)
	userGroup.DELETE("", userHandler.UserDeleteById)

	log.Fatal(e.Start(":8083"))

}
