package main

import (
	"MessagesService/databases"
	"MessagesService/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)


const _APIVersion = "v1"
const _prefix = "/api/" + _APIVersion

func main() {
	godotenv.Load()

	databases.InitDatabases()

	e := echo.New()

	h := &handlers.Handler{}

	e.POST(_prefix + "/channels/:channelID/messages", h.NewMessageHandler)
	e.GET(_prefix + "/channels/:channelID/messages", h.GetMessagesHandler)

	e.Start(":1323")
}