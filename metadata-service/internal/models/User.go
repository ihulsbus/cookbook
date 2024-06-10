package models

type User struct {
	Provider string `json:"provider"`
	UserID   string `json:"user_id"`
}
