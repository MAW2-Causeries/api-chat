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
var getChannelUsers = models.GetChannelUsers
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

	for {
		_, _, err := readConnectionMessage(ws)
		if err != nil {
			break
		}
	}

	delete(connectedUser, userID)
	return nil
}

func (h *Handler) notify(channelID string, message *models.Message) {
	for _, userID := range getChannelUsers(channelID) {
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
