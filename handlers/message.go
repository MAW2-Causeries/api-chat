package handlers

import (
	"MessagesService/models"

	"github.com/labstack/echo/v4"
)

// NewMessageHandler handles the /api/v1/message route for Post requests
func (h *Handler) NewMessageHandler(c echo.Context) (err error) {
	if c.Request().Method != echo.POST {
		return echo.NewHTTPError(405, "Method Not Allowed")
	}

	authorID := c.FormValue("author_id")
	channelID := c.FormValue("channel_id")
	content := c.FormValue("content")

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