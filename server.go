package main

import (
	"MessagesService/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)


const _APIVersion = "v1"
const _prefix = "/api/" + _APIVersion

func main() {
	e := echo.New()
	e.GET(_prefix, _RootHandler)

	h := &handlers.Handler{}

	e.GET(_prefix + "/messages", h.MessageHandler)

	e.Start(":1323")
}

func _RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}