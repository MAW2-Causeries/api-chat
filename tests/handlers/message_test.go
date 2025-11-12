package tests

import (
	"MessagesService/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMessageHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handlers.Handler{}

	if assert.NoError(t, h.MessageHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "This is the message endpoint!", rec.Body.String())
	}
}