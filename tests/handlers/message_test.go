package tests

import (
	"MessagesService/handlers"
	"MessagesService/models"
	"MessagesService/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandlerReturnNewMessage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("author_id=DD79A816-D227-4AEC-B413-6AF520B0B157&content=hello+world"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return &models.Message{
			ID:       	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: 	"f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   	"feudbfuidsfhdosr",
		}
	})
	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", nil
	})

	h := &handlers.Handler{}

	if assert.NoError(t, h.NewMessageHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", resp["author_id"])
			assert.Equal(t, "f63f7c42-c567-4b17-bd3a-93c1eb510ed9", resp["channel_id"])
			assert.Equal(t, "feudbfuidsfhdosr", resp["content"])
		}
	}
}

func TestNewMessageHandlerMissingFields(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", nil
	})
	
	h := &handlers.Handler{}
	
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing required fields", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessagesHandlerReturnMessages(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", nil
	})
	
	monkey.Patch(models.GetMessagesByChannelID, func(channelID string, limit, page int) ([]*models.Message, error) {
		return []*models.Message{
			{
				ID:       	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
				AuthorID:  	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
				ChannelID: 	"f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
				Content:   	"feudbfuidsfhdosr",
			},
		}, nil
	})
	
	h := &handlers.Handler{}
	if assert.NoError(t, h.GetMessagesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		
		var resp []map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Len(t, resp, 1)
			assert.Equal(t, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", resp[0]["author_id"])
			assert.Equal(t, "f63f7c42-c567-4b17-bd3a-93c1eb510ed9", resp[0]["channel_id"])
			assert.Equal(t, "feudbfuidsfhdosr", resp[0]["content"])
		}
	}
}

func TestGetMessagesHandlerMissingChannelID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels//messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("")

	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", nil
	})
	
	h := &handlers.Handler{}
	
	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing channel ID", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessageUnauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "", errors.New("invalid token")
	})
	
	h := &handlers.Handler{}
	
	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusUnauthorized, he.Code)
		assert.Equal(t, "invalid token", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestNewMessageWithJsonBody(t *testing.T) {
	e := echo.New()
	jsonBody := `{"content":"Hello, JSON!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return &models.Message{
			ID:       	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  	"bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: 	"f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   	"Hello, JSON!",
		}
	})
	monkey.Patch(utils.VerifyBearerToken, func(Authorization string) (string, error) {
		return "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", nil
	})

	h := &handlers.Handler{}

	if assert.NoError(t, h.NewMessageHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", resp["author_id"])
			assert.Equal(t, "f63f7c42-c567-4b17-bd3a-93c1eb510ed9", resp["channel_id"])
			assert.Equal(t, "Hello, JSON!", resp["content"])
		}
	}
}