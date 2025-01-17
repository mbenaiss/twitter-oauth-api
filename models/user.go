package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"myapp/auth"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserResponse struct {
	Data User `json:"data"`
}

// GetUserInfo fetches the user information from Twitter API
func GetUserInfo(accessToken string) (*User, error) {
	req, err := http.NewRequest("GET", auth.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating user info request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info endpoint returned status %d", resp.StatusCode)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("decoding user info response: %w", err)
	}

	return &userResp.Data, nil
}
