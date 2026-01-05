package handlers

import (
	"MessagesService/models"
	"MessagesService/utils"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connected_user = make(map[string]*websocket.Conn)

func (h *Handler) Websocket(c echo.Context) (err error) {
	authHeader := c.Request().Header.Get("Authorization")
	authorID, err := utils.VerifyBearerToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	connected_user[authorID] = ws

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	connected_user[authorID] = nil
	return nil
}

func (h *Handler) notify(channelID string, message *models.Message) {
	for _, conn := range connected_user {
		err := conn.WriteJSON(message.ToMap())
		if err != nil {
			continue
		}
	}
}

// NewMessageHandler handles route for Post requests
func (h *Handler) NewMessageHandler(c echo.Context) (err error) {
	channelID := c.Param("channelID")
	content := ""
	if c.Request().Header.Get("Content-Type") == "application/json" {
		body := struct {
			Content string `json:"content"`
		}{}
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(400, "Invalid JSON body")
		}
		content = body.Content
	} else {
		content = c.FormValue("content")
	}
	authorID := ""

	if channelID == "" || content == "" {
		return echo.NewHTTPError(400, "Missing required fields")
	}

	authHeader := c.Request().Header.Get("Authorization")
	authorID, err = utils.VerifyBearerToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	message := models.NewMessage(
		authorID,
		channelID,
		content,
	)

	if message == nil {
		return echo.NewHTTPError(500, "Failed to create message")
	}

	go h.notify(channelID, message)

	return c.JSON(201, message.ToMap())
}

// GetMessagesHandler handles the /api/v1/message route for Get requests
func (h *Handler) GetMessagesHandler(c echo.Context) (err error) {
	channelID := c.Param("channelID")
	limit := c.QueryParam("limit")
	page := c.QueryParam("page")

	authHeader := c.Request().Header.Get("Authorization")
	_, err = utils.VerifyBearerToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	if limit == "" {
		limit = "50"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		return echo.NewHTTPError(400, "Invalid limit; must be an integer between 1 and 100")
	}

	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		return echo.NewHTTPError(400, "Invalid page; must be an integer greater than 0")
	}
	
	if channelID == "" {
		return echo.NewHTTPError(400, "Missing channel ID")
	}

	messages, err := models.GetMessagesByChannelID(channelID, limitInt, pageInt)
	if err != nil {
		println(err.Error())
		return echo.NewHTTPError(500, "Failed to retrieve messages")
	}

	messageMaps := make([]map[string]any, len(messages))
	for i, message := range messages {
		messageMaps[i] = message.ToMap()
	}

	return c.JSON(200, messageMaps)
}

func (h *Handler) GetMessageHandler(c echo.Context) (err error) {
	channelID := c.Param("channelID")
	messageID := c.Param("messageID")

	authHeader := c.Request().Header.Get("Authorization")
	_, err = utils.VerifyBearerToken(authHeader)
	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}
	
	if channelID == "" || messageID == "" {
		return echo.NewHTTPError(400, "Missing required fields")
	}
	message := models.GetMessageByChannelIdAndMessageID(channelID, messageID)
	if message == nil {
		return echo.NewHTTPError(404, "Message not found")
	}

	return c.JSON(200, message.ToMap())
}