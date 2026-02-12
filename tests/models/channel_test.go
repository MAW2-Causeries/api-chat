package tests

import (
	"MessagesService/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestGetUserChannels(t *testing.T) {
	fakeUserId := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fakeResponseBody := `["27731CCA-ADB5-42DB-AA8C-500994FC4098","3F2504E0-4F89-11D3-9A0C-0305E82C3301"]`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fakeResponseBody))
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	usersChannels := models.GetUserChannels(fakeUserId)

	expectedChannels := []string{
		"27731CCA-ADB5-42DB-AA8C-500994FC4098",
		"3F2504E0-4F89-11D3-9A0C-0305E82C3301",
	}

	assert.Equal(t, expectedChannels, usersChannels)
}

func TestDoesUserCanSendMessageInChannel(t *testing.T) {
	fakeUserId := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelId := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	canSend := models.DoesUserCanSendMessageInChannel(fakeUserId, fakeChannelId)

	assert.True(t, canSend)
}

func TestDoesUserCanReadMessagesInChannel(t *testing.T) {
	fakeUserId := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelId := "27731CCA-ADB5-42DB-AA8C-500994FC4098"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("BASE_API_URL", server.URL+"/api/v1")

	canRead := models.DoesUserCanReadMessagesInChannel(fakeUserId, fakeChannelId)

	assert.True(t, canRead)
}