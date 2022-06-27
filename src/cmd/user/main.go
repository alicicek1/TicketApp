package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"ticketApp/src/config"
	"ticketApp/src/handler"
	"ticketApp/src/repository"
	"ticketApp/src/service"
)

func main() {
	mCfg := config.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()
	//e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

	userCollection := mCfg.GetCollection(client, cfg.UserColName)
	userRepository := repository.NewUserRepository(userCollection)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, cfg)
	userGroup := e.Group("/api/users")
	//userGroup.GET("", userHandler.UserGetById)
	userGroup.GET("/:id", userHandler.UserGetById)
	userGroup.GET("", userHandler.UserGetAll)
	userGroup.POST("", userHandler.UserUpsert)
	userGroup.POST("/login", userHandler.Login)
	userGroup.DELETE("/:id", userHandler.UserDeleteById)

	log.Fatal(e.Start(":8083"))

}
