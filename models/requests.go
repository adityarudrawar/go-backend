package models

type PostMessageRequest struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}