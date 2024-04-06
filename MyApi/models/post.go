package models

import (
	"time"
)

type Post struct {
	ID		  string	`json:"_id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UserID    string    `json:"userId" bson:"userId"`
	FirstName string    `json:"firstName" bson:"firstName"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Comments  []Comment `json:"comments,omitempty" bson:"comments,omitempty"`
	UpVotes   []string  `json:"upVotes,omitempty" bson:"upVotes,omitempty"`
}

type Comment struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	ID        string    `json:"id" bson:"id"`
	FirstName string    `json:"firstName" bson:"firstName"`
	Content   string    `json:"content" bson:"content"`
}