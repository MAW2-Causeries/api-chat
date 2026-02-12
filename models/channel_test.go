package models

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)


func TestGetUserChannels(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fakeResponseBody := `["27731CCA-ADB5-42DB-AA8C-500994FC4098","3F2504E0-4F89-11D3-9A0C-0305E82C3301"]`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fakeResponseBody))
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	usersChannels := GetUserChannels(fakeUserID)

	expectedChannels := []string{
		"27731CCA-ADB5-42DB-AA8C-500994FC4098",
		"3F2504E0-4F89-11D3-9A0C-0305E82C3301",
	}

	assert.Equal(t, expectedChannels, usersChannels)
}

func TestGetUserChannelsWithError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"

	monkey.Patch(http.Get, func(url string) (*http.Response, error) {
		return nil, assert.AnError
	})

	defer monkey.Unpatch(http.Get)


	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	usersChannels := GetUserChannels(fakeUserID)
	assert.Empty(t, usersChannels)
}

func TestGetUserChannelsWithIOBodyError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	monkey.Patch(io.ReadAll, func(r io.Reader) ([]byte, error) {
		return nil, assert.AnError
	})

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	usersChannels := GetUserChannels(fakeUserID)
	assert.Empty(t, usersChannels)
}

func TestDoesUserCanSendMessageInChannel(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)

	assert.True(t, canSend)
}

func TestDoesUserCanSendMessageInChannelWithError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	monkey.Patch(http.Get, func(url string) (*http.Response, error) {
		return nil, assert.AnError
	})
	defer monkey.Unpatch(http.Get)
	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)
	assert.False(t, canSend)
}

func TestDoesUserCanSendMessageInChannelWithNonOKStatus(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)
	assert.False(t, canSend)
}

func TestDoesUserCanReadMessagesInChannel(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	canRead := DoesUserCanReadMessagesInChannel(fakeUserID, fakeChannelID)

	assert.True(t, canRead)
}