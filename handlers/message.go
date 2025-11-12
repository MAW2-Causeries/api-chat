package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// MessageHandler handles the /api/v1/message route
func (h *Handler) MessageHandler(c echo.Context) (err error) {
	return c.String(http.StatusOK, "This is the message endpoint!")
}