package handlers

import (
	"MessagesService/models"
	"MessagesService/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

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

// GetMessageHandler handles the /api/v1/message/:messageID route for Get requests
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
	message := models.GetMessageByChannelIDAndMessageID(channelID, messageID)
	if message == nil {
		return echo.NewHTTPError(404, "Message not found")
	}

	return c.JSON(200, message.ToMap())
}