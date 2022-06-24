package main

import (
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
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()
	e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

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

	categoryCollection := mCfg.GetCollection(client, cfg.CategoryColName)
	categoryRepository := repository.NewCategoryRepository(categoryCollection)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService, cfg)
	categoryGroup := e.Group("/api/categories")
	categoryGroup.GET("/:id", categoryHandler.CategoryGetById)
	categoryGroup.GET("", categoryHandler.CategoryGetAll)
	categoryGroup.POST("", categoryHandler.CategoryInsert)
	categoryGroup.DELETE("/:id", categoryHandler.CategoryDeleteById)

	ticketCollection := mCfg.GetCollection(client, cfg.TicketColName)
	ticketRepository := repository.NewTicketRepository(ticketCollection)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)
	ticketGroup := e.Group("/api/tickets")
	ticketGroup.GET("/:id", ticketHandler.TicketGetById)
	ticketGroup.GET("", ticketHandler.TicketGetAll)
	ticketGroup.POST("", ticketHandler.TicketInsert)
	ticketGroup.DELETE("/:id", ticketHandler.TicketDeleteById)

	log.Fatal(e.Start(":8083"))

}
