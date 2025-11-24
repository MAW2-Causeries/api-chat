package main

import (
	"MessagesService/databases"
	"MessagesService/handlers"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)


const _APIVersion = "v1"
const _prefix = "/api/" + _APIVersion

func main() {
	godotenv.Load()

	databases.InitDatabases()

	e := echo.New()
	e.GET(_prefix, _RootHandler)

	h := &handlers.Handler{}

	e.POST(_prefix + "/channels/:channelID/messages", h.NewMessageHandler)
	e.GET(_prefix + "/channels/:channelID/messages", h.GetMessagesHandler)

	e.Start(":1323")
}

func _RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}