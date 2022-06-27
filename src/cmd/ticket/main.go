package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"ticketApp/src/config"
	"ticketApp/src/handler"
	"ticketApp/src/repository"
	"ticketApp/src/service"
)

func zeroVal(i int) {
	i = 0
}

func zeroPtr(i *int) {
	*i = 0
}

func main() {
	mCfg := config.NewMongoConfig()
	client, _, cancel, cfg := mCfg.ConnectDatabase()
	defer cancel()

	e := echo.New()
	//e.HTTPErrorHandler = util.NewHttpErrorHandler(util.NewErrorStatusCodeMaps()).Handler

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
