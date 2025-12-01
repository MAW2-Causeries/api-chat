package handlers

import (
	"MessagesService/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// NewMessageHandler handles the /api/v1/message route for Post requests
func (h *Handler) NewMessageHandler(c echo.Context) (err error) {
	channelID := c.Param("channelID")
	content := c.FormValue("content")
	authorID := c.FormValue("author_id")

	if authorID == "" || channelID == "" || content == "" {
		return echo.NewHTTPError(400, "Missing required fields")
	}

	message := models.NewMessage(
		authorID,
		channelID,
		content,
	)

	if message == nil {
		return echo.NewHTTPError(500, "Failed to create message")
	}

	return c.JSON(201, message.ToMap())
}

// GetMessagesHandler handles the /api/v1/message route for Get requests
func (h *Handler) GetMessagesHandler(c echo.Context) (err error) {
	channelID := c.Param("channelID")
	limit := c.QueryParam("limit")

	if limit == "" {
		limit = "50"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		return echo.NewHTTPError(400, "Invalid limit; must be an integer between 1 and 100")
	}
	
	if channelID == "" {
		return echo.NewHTTPError(400, "Missing channel ID")
	}

	messages, err := models.GetMessagesByChannelID(channelID, limitInt)
	if err != nil {
		println(err.Error())
		return echo.NewHTTPError(500, "Failed to retrieve messages")
	}

	messageMaps := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		messageMaps[i] = message.ToMap()
	}

	return c.JSON(200, messageMaps)
}