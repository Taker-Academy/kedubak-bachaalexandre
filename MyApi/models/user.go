package models

import (
    "time"
)

type User struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Password    string    `json:"password"`
	LastUpVote  time.Time `json:"lastUpVote"`
}