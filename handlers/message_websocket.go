package handlers

import (
	"MessagesService/models"
	"MessagesService/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connectedUser = make(map[string]*websocket.Conn)
var channelsSubscribtion = make(map[string][]string)

// Websocket handles websocket connections for real-time message updates
func (h *Handler) Websocket(c echo.Context) (err error) {
	authHeader := c.Request().Header.Get("Authorization")
	userID, err := utils.VerifyBearerToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	connectedUser[userID] = ws

	for _, channelID := range models.GetUserChannels(userID) {
		channelsSubscribtion[channelID] = append(channelsSubscribtion[channelID], userID)
	}

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	connectedUser[userID] = nil
	for _, channelID := range models.GetUserChannels(userID) {
		subscribers := channelsSubscribtion[channelID]
		for i, id := range subscribers {
			if id == userID {
				channelsSubscribtion[channelID] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (h *Handler) notify(channelID string, message *models.Message) {
	for _, userID := range channelsSubscribtion[channelID] {
		conn := connectedUser[userID]
		if conn == nil {
			continue
		}

		err := conn.WriteJSON(message.ToMap())
		if err != nil {
			continue
		}
	}
}
