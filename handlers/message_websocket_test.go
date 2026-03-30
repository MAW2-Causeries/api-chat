package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cpnv.ch/messagesservice/middlewares"
	"cpnv.ch/messagesservice/models"
	"github.com/bouk/monkey"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessageHandlerFormBodySuccess(t *testing.T) {
	defer monkey.UnpatchAll()

	e := echo.New()
	formBody := "content=hello+form"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/channels/channel-1/messages", strings.NewReader(formBody))
	req = req.WithContext(context.WithValue(context.Background(), middlewares.UserIDKey, "user-1"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("channel-1")

	monkey.Patch(models.DoesUserCanSendMessageInChannel, func(userID, channelID string) bool {
		return true
	})
	monkey.Patch(models.NewMessage, func(authorID, channelID, content string) *models.Message {
		return &models.Message{ID: "1", AuthorID: authorID, ChannelID: channelID, Content: content}
	})

	err := (&Handler{}).NewMessageHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "\"content\":\"hello form\"")
}

func TestGetMessagesHandlerLimitTooLarge(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/channels/channel-1/messages?limit=101", nil)
	req = req.WithContext(context.WithValue(context.Background(), middlewares.UserIDKey, "user-1"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("channelID")
	c.SetParamValues("channel-1")

	err := (&Handler{}).GetMessagesHandler(c)
	he, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, he.Code)
	assert.Equal(t, "Invalid limit; must be an integer between 1 and 100", he.Message)
}

func TestWebsocketReturnsUpgradeError(t *testing.T) {
	oldUpgradeConnection := upgradeConnection
	t.Cleanup(func() {
		upgradeConnection = oldUpgradeConnection
	})
	upgradeConnection = func(c echo.Context) (*websocket.Conn, error) {
		return nil, errors.New("upgrade failed")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req = req.WithContext(context.WithValue(context.Background(), middlewares.UserIDKey, "user-1"))
	rec := httptest.NewRecorder()

	err := (&Handler{}).Websocket(e.NewContext(req, rec))

	assert.Error(t, err)
}

func TestWebsocketSubscribesAndUnsubscribesUser(t *testing.T) {
	connectedUser = make(map[string]*websocket.Conn)
	oldUpgradeConnection := upgradeConnection
	oldCloseConnection := closeConnection
	oldReadConnectionMessage := readConnectionMessage
	t.Cleanup(func() {
		upgradeConnection = oldUpgradeConnection
		closeConnection = oldCloseConnection
		readConnectionMessage = oldReadConnectionMessage
	})

	readCalls := 0
	upgradeConnection = func(c echo.Context) (*websocket.Conn, error) {
		return &websocket.Conn{}, nil
	}
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	readConnectionMessage = func(conn *websocket.Conn) (int, []byte, error) {
		readCalls++
		if readCalls == 1 {
			assert.NotNil(t, connectedUser["user-1"])
		}
		return 0, nil, errors.New("closed")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req = req.WithContext(context.WithValue(context.Background(), middlewares.UserIDKey, "user-1"))
	rec := httptest.NewRecorder()

	err := (&Handler{}).Websocket(e.NewContext(req, rec))

	require.NoError(t, err)
	assert.Nil(t, connectedUser["user-1"])
	assert.Len(t, connectedUser, 0)
}

func TestNotifyUsesChannelUsersAndSkipsDisconnectedUsers(t *testing.T) {
	connectedUser = map[string]*websocket.Conn{
		"nil-user": nil,
		"ok-user":  {},
		"bad-user": {},
	}
	oldGetChannelUsers := getChannelUsers
	oldWriteConnectionJSON := writeConnectionJSON
	t.Cleanup(func() {
		getChannelUsers = oldGetChannelUsers
		writeConnectionJSON = oldWriteConnectionJSON
	})

	getChannelUsers = func(channelID string) []string {
		assert.Equal(t, "channel-1", channelID)
		return []string{"nil-user", "ok-user", "bad-user"}
	}

	var writes int
	writeConnectionJSON = func(conn *websocket.Conn, v any) error {
		writes++
		if conn == connectedUser["bad-user"] {
			return assert.AnError
		}
		assert.Equal(t, "hello", v.(map[string]any)["content"])
		return nil
	}

	(&Handler{}).notify("channel-1", &models.Message{Content: "hello"})

	assert.Equal(t, 2, writes)
}
