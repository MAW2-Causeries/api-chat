package handlers

import (
	"cpnv.ch/messagesservice/middlewares"
	"cpnv.ch/messagesservice/models"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connectedUser = make(map[string]*websocket.Conn)
var channelsSubscribtion = make(map[string][]string)
var getUserChannels = models.GetUserChannels
var upgradeConnection = func(c echo.Context) (*websocket.Conn, error) {
	return upgrader.Upgrade(c.Response(), c.Request(), nil)
}
var closeConnection = func(conn *websocket.Conn) error {
	return conn.Close()
}
var readConnectionMessage = func(conn *websocket.Conn) (int, []byte, error) {
	return conn.ReadMessage()
}
var writeConnectionJSON = func(conn *websocket.Conn, value any) error {
	return conn.WriteJSON(value)
}

// Websocket handles websocket connections for real-time message updates
func (h *Handler) Websocket(c echo.Context) (err error) {
	userID := c.Request().Context().Value(middlewares.UserIDKey).(string)

	ws, err := upgradeConnection(c)
	if err != nil {
		return err
	}
	defer closeConnection(ws)

	connectedUser[userID] = ws

	for _, channelID := range getUserChannels(userID) {
		channelsSubscribtion[channelID] = append(channelsSubscribtion[channelID], userID)
	}

	for {
		_, _, err := readConnectionMessage(ws)
		if err != nil {
			break
		}
	}

	connectedUser[userID] = nil
	for _, channelID := range getUserChannels(userID) {
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

		err := writeConnectionJSON(conn, message.ToMap())
		if err != nil {
			continue
		}
	}
}
