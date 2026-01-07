package models

import (
	"MessagesService/utils"
	"encoding/json"
	"io"
	"net/http"
)

func GetUserChannels(userID string) []string {
	var channelIDs []string

	baseURL := utils.GetEnv("BASE_API_URL", "http://localhost:8080/api/v1")
	resp, err := http.Get(baseURL + "/users/" + userID + "/channels?field=id")
	if err != nil {
		return channelIDs
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return channelIDs
	}

	json.Unmarshal(body, &channelIDs)
	return channelIDs
}

func DoesUserCanSendMessageInChannel(userID, channelID string) bool {
	baseURL := utils.GetEnv("BASE_API_URL", "http://localhost:8080/api/v1")
	resp, err := http.Get(baseURL + "/channels/" + channelID + "/users/" + userID)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func DoesUserCanReadMessagesInChannel(userID, channelID string) bool {
	return DoesUserCanSendMessageInChannel(userID, channelID)
}