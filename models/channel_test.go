package models

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserChannels(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	oldDoHTTPRequest := doHTTPRequest
	oldReadHTTPBody := readHTTPBody
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
		readHTTPBody = oldReadHTTPBody
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "http://localhost:8080/api/v1/users/"+fakeUserID+"/channels?field=id", req.URL.String())
		assert.Equal(t, "super-secret-token", req.Header.Get("X-Master-Secret-Token"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`["27731CCA-ADB5-42DB-AA8C-500994FC4098","3F2504E0-4F89-11D3-9A0C-0305E82C3301"]`)),
		}, nil
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")
	t.Setenv("MASTER_SECRET_TOKEN", "super-secret-token")

	usersChannels := GetUserChannels(fakeUserID)

	expectedChannels := []string{
		"27731CCA-ADB5-42DB-AA8C-500994FC4098",
		"3F2504E0-4F89-11D3-9A0C-0305E82C3301",
	}

	assert.Equal(t, expectedChannels, usersChannels)
}

func TestGetUserChannelsWithError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
	oldDoHTTPRequest := doHTTPRequest
	t.Cleanup(func() {
		doHTTPRequest = oldDoHTTPRequest
	})

	doHTTPRequest = func(req *http.Request) (*http.Response, error) {
		return nil, assert.AnError
	}

	t.Setenv("BASE_API_URL", "http://localhost:8080/api/v1")

	usersChannels := GetUserChannels(fakeUserID)
	assert.Empty(t, usersChannels)
}

func TestGetUserChannelsWithIOBodyError(t *testing.T) {
	fakeUserID := "F77AC4EA-4AF0-4F64-A985-CAA0284C8257"
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

	usersChannels := GetUserChannels(fakeUserID)
	assert.Empty(t, usersChannels)
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
