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

	err := databases.InitDatabases()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	h := &handlers.Handler{}

	e.POST(_prefix + "/channels/:channelID/messages", h.NewMessageHandler)
	e.GET(_prefix + "/channels/:channelID/messages", h.GetMessagesHandler)
	e.GET(_prefix + "/channels/:channelID/messages/:messageID", h.GetMessageHandler)
	e.GET(_prefix, h.Websocket)

	e.Start(":1323")
}