package main

import (
	"github.com/labstack/echo/v4"
	"nat.vg/narwhal/internal/handler"
	"nat.vg/narwhal/internal/service"
)

func main() {
	e := echo.New()
	repo := service.NewNarwhalService()

	handler.NarwhalHandler(e, repo)

	e.Logger.Fatal(e.Start(":8080"))
}
