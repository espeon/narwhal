package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"nat.vg/narwhal/internal/handler"
	"nat.vg/narwhal/internal/service"
)

func main() {
	godotenv.Load("narwhal.env", ".env")
	e := echo.New()
	repo := service.NewNarwhalService()

	handler.NarwhalHandler(e, repo)

	// get port from env
	port := os.Getenv("NARWHAL_PORT")
	if port == "" {
		port = "46449"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
