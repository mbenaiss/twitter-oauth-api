package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbenaiss/twitter-oauth-api/models"
)

const (
	userInfoURL = "https://api.x.com/2/users/me"
)

type userResponse struct {
	Data models.User `json:"data"`
}

// GetUserInfo fetches the user information from Twitter API
func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (models.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("creating user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.User{}, fmt.Errorf("executing user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.User{}, fmt.Errorf("user info endpoint returned status %d", resp.StatusCode)
	}

	var userResp userResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return models.User{}, fmt.Errorf("decoding user info response: %w", err)
	}

	return userResp.Data, nil
}
