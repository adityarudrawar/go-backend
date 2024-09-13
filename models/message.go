package models
import "time"

type Message struct {
	Uuid      string `json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	Sender    string `json:"sender"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
	Content   string `json:"content"`
}