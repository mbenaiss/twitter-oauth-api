package models

// User represents the user information returned by the Twitter API
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
