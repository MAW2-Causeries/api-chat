package models

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChannelUsers(t *testing.T) {
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	oldReadHTTPBody := readHTTPBody
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
		readHTTPBody = oldReadHTTPBody
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "http://localhost:8080/api/v1/channels/"+fakeChannelID+"/users", req.URL.String())
		assert.Equal(t, "super-secret-token", req.Header.Get("X-Master-Secret-Token"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`[{"id":"F77AC4EA-4AF0-4F64-A985-CAA0284C8257"},{"id":"3F2504E0-4F89-11D3-9A0C-0305E82C3301"}]`)),
		}, nil
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")
	t.Setenv("MASTER_SECRET_TOKEN", "super-secret-token")

	channelUsers := GetChannelUsers(fakeChannelID)

	expectedUsers := []string{
		"F77AC4EA-4AF0-4F64-A985-CAA0284C8257",
		"3F2504E0-4F89-11D3-9A0C-0305E82C3301",
	}

	assert.Equal(t, expectedUsers, channelUsers)
}

func TestGetChannelUsersWithError(t *testing.T) {
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return nil, assert.AnError
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	channelUsers := GetChannelUsers(fakeChannelID)
	assert.Empty(t, channelUsers)
}

func TestGetChannelUsersWithIOBodyError(t *testing.T) {
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	oldReadHTTPBody := readHTTPBody
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
		readHTTPBody = oldReadHTTPBody
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ignored")),
		}, nil
	}
	readHTTPBody = func(r io.Reader) ([]byte, error) {
		return nil, assert.AnError
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	channelUsers := GetChannelUsers(fakeChannelID)
	assert.Empty(t, channelUsers)
}

func TestDoesUserCanSendMessageInChannel(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "http://localhost:8080/api/v1/channels/"+fakeChannelID+"/users/"+fakeUserID, req.URL.String())
		assert.Equal(t, "super-secret-token", req.Header.Get("X-Master-Secret-Token"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")
	t.Setenv("MASTER_SECRET_TOKEN", "super-secret-token")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)

	assert.True(t, canSend)
}

func TestDoesUserCanSendMessageInChannelWithError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return nil, assert.AnError
	}
	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)
	assert.False(t, canSend)
}

func TestDoesUserCanSendMessageInChannelWithNonOKStatus(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	}
	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	canSend := DoesUserCanSendMessageInChannel(fakeUserID, fakeChannelID)
	assert.False(t, canSend)
}

func TestDoesUserCanReadMessagesInChannel(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	fakeChannelID := "27731CCA-ADB5-42DB-AA8C-500994FC4098"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	canRead := DoesUserCanReadMessagesInChannel(fakeUserID, fakeChannelID)

	assert.True(t, canRead)
}
