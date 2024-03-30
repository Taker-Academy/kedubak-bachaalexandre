package models

import (
	"time"
)

type Comment struct {
    CreatedAt time.Time `json:"createdAt"`
    ID        string    `json:"id"`
    FirstName string    `json:"firstName"`
    Content   string    `json:"content"`
}

type Post struct {
    CreatedAt time.Time `json:"createdAt"`
    UserID    string    `json:"userId"`
    FirstName string    `json:"firstName"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Comments  []Comment `json:"comments"`
    UpVotes   []string  `json:"upVotes"`
}