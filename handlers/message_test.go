package handlers

import (
	"context"
	"encoding/json"
	"cpnv.ch/messagesservice/middlewares"
	"cpnv.ch/messagesservice/models"
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
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("{\"content\":\"hello world\"}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return &models.Message{
			ID:        "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: "f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   "hello world",
		}
	})

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	h := &Handler{}

	if assert.NoError(t, h.NewMessageHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", resp["author_id"])
			assert.Equal(t, "f63f7c42-c567-4b17-bd3a-93c1eb510ed9", resp["channel_id"])
			assert.Equal(t, "hello world", resp["content"])
		}
	}
}

func TestNewMessageHandlerInvalidJSONBody(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("{invalid json}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Invalid JSON body", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestNewMessageHandlerMissingFields(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader(""))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}

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
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	monkey.Patch(models.GetMessagesByChannelID, func(channelID string, limit, page int) []*models.Message {
		return []*models.Message{
			{
				ID:        "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
				AuthorID:  "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
				ChannelID: "f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
				Content:   "feudbfuidsfhdosr",
			},
		}
	})

	h := &Handler{}
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
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels//messages", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("")

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing channel ID", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestNewMessageWithJsonBody(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	jsonBody := `{"content":"Hello, JSON!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return &models.Message{
			ID:        "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: "f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   "Hello, JSON!",
		}
	})

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	h := &Handler{}

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

func TestGetMessageHandlerReturnMessage(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages/BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID", "messageID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097", "BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	monkey.Patch(models.GetMessageByChannelIDAndMessageID, func(channelID string, messageID string) *models.Message {
		return &models.Message{
			ID:        "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: "f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   "feudbfuidsfhdosr",
		}
	})

	h := &Handler{}
	if assert.NoError(t, h.GetMessageHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70", resp["author_id"])
			assert.Equal(t, "f63f7c42-c567-4b17-bd3a-93c1eb510ed9", resp["channel_id"])
			assert.Equal(t, "feudbfuidsfhdosr", resp["content"])
		}
	}
}

func TestGetMessageHandlerMessageNotFound(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages/BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID", "messageID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097", "BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70")

	monkey.Patch(models.GetMessageByChannelIDAndMessageID, func(channelID string, messageID string) *models.Message {
		return nil
	})

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	h := &Handler{}

	err := h.GetMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusNotFound, he.Code)
		assert.Equal(t, "Message not found", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessageHandlerInvalidLimit(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages?limit=invalid", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Invalid limit; must be an integer between 1 and 100", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}
func TestGetMessagesHandlerInvalidPage(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages?page=invalid", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Invalid page; must be an integer greater than 0", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessagesHandlerNegativePage(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages?page=-1", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Invalid page; must be an integer greater than 0", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessagesHandlerZeroPage(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages?page=0", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Invalid page; must be an integer greater than 0", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}
func TestGetMessagesHandlerUserNoPermission(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return false
	})

	h := &Handler{}

	err := h.GetMessagesHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusForbidden, he.Code)
		assert.Equal(t, "User does not have permission to read messages in this channel", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessageHandlerUserNoPermission(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages/BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID", "messageID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097", "BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return false
	})

	monkey.Patch(models.GetMessageByChannelIDAndMessageID, func(channelID string, messageID string) *models.Message {
		return &models.Message{
			ID:        "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			AuthorID:  "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70",
			ChannelID: "f63f7c42-c567-4b17-bd3a-93c1eb510ed9",
			Content:   "feudbfuidsfhdosr",
		}
	})

	h := &Handler{}

	err := h.GetMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusForbidden, he.Code)
		assert.Equal(t, "User does not have permission to read messages in this channel", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}
func TestGetMessagesHandlerNilMessagesReturnsEmptyList(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	monkey.Patch(models.GetMessagesByChannelID, func(channelID string, limit, page int) []*models.Message {
		return nil
	})

	h := &Handler{}

	if assert.NoError(t, h.GetMessagesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []map[string]any
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Empty(t, resp)
		}
	}
}
func TestNewMessageHandlerMissingChannelID(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels//messages", strings.NewReader("{\"content\":\"hello world\"}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("")

	h := &Handler{}
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing required fields", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestNewMessageHandlerMissingContent(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("{\"content\":\"\"}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	h := &Handler{}
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing required fields", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessageHandlerMissingChannelID(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels//messages/BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID", "messageID")
	c.SetParamValues("", "BB6A2B8A-954A-4AC2-A7B9-4B5A100AFB70")

	h := &Handler{}
	err := h.GetMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing required fields", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestGetMessageHandlerMissingMessageID(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages/", nil)
	req = req.WithContext(requestContext)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID", "messageID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097", "")

	h := &Handler{}
	err := h.GetMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, "Missing required fields", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}
func TestNewMessageHandlerUserNoPermission(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("{\"content\":\"hello world\"}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return false
	})

	h := &Handler{}
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusForbidden, he.Code)
		assert.Equal(t, "User does not have permission to send messages in this channel", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}

func TestNewMessageHandlerFailedToCreateMessage(t *testing.T) {
	e := echo.New()
	requestContext := context.WithValue(context.Background(), middlewares.UserIDKey, "bb6a2b8a-954a-4ac2-a7b9-4b5a100afb70")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/DD04A392-A4D6-45F5-86B5-E070E7588097/messages", strings.NewReader("{\"content\":\"hello world\"}"))
	req = req.WithContext(requestContext)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("DD04A392-A4D6-45F5-86B5-E070E7588097")

	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return nil
	})

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})

	h := &Handler{}
	err := h.NewMessageHandler(c)
	if he, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusInternalServerError, he.Code)
		assert.Equal(t, "Failed to create message", he.Message)
	} else {
		t.Fatalf("expected HTTPError, got %v", err)
	}
}
