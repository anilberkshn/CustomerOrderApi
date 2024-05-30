package main

import (
	"CustomerOrderApi/Customer/handlers"
	"CustomerOrderApi/Customer/repositories"
	"CustomerOrderApi/Customer/services"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	repository := repositories.NewRepository()
	service := services.NewService(repository)
	handler := handlers.NewHandler(service, e, repository)

	handler.InitEndpoints()

	e.Logger.Fatal(e.Start(":8080"))
}
