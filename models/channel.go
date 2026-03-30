package models

import (
	"encoding/json"
	"io"
	"net/http"

	"cpnv.ch/messagesservice/utils"
)

var doHTTPRequest = http.DefaultClient.Do
var readHTTPBody = io.ReadAll

func getHTTP(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Master-Secret-Token", utils.GetEnv("MASTER_SECRET_TOKEN", ""))

	return doHTTPRequest(req)
}

// GetChannelUsers retrieves the list of user IDs that are members of a specific channel by making an HTTP GET request to the API. It returns a slice of user IDs, or an empty slice if there was an error during the request or response parsing.
func GetChannelUsers(channelID string) []string {
	var users []struct {
		ID string `json:"id"`
	}

	baseURL := utils.GetEnv("BASE_API_URL", "http://localhost:8080/api/v1")
	resp, err := getHTTP(baseURL + "/channels/" + channelID + "/users")
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, err := readHTTPBody(resp.Body)
	if err != nil {
		return []string{}
	}

	if err := json.Unmarshal(body, &users); err != nil {
		return []string{}
	}

	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}
	return userIDs
}

// DoesUserCanSendMessageInChannel checks if a user has permission to send messages in a specific channel by making an HTTP GET request to the API. It returns true if the user can send messages, and false otherwise.
func DoesUserCanSendMessageInChannel(userID, channelID string) bool {
	baseURL := utils.GetEnv("BASE_API_URL", "http://localhost:8080/api/v1")
	resp, err := getHTTP(baseURL + "/channels/" + channelID + "/users/" + userID)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

// DoesUserCanReadMessagesInChannel checks if a user has permission to read messages in a specific channel by delegating to DoesUserCanSendMessageInChannel. It returns true if the user can read messages, and false otherwise. This function may be modified in the future if the permissions for reading and sending messages differ.
// TODO: This function is currently the same as DoesUserCanSendMessageInChannel, but it can be modified in the future if the permissions for reading and sending messages differ.
func DoesUserCanReadMessagesInChannel(userID, channelID string) bool {
	return DoesUserCanSendMessageInChannel(userID, channelID)
}
