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
