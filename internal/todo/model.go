package todo

import (
	"time"
)

// Todo model
type Todo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedOn time.Time `json:"created_on"`
}
