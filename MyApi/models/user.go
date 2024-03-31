package models

import (
    "time"
)

type User struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	Email       string    `json:"email" bson:"email"`
	FirstName   string    `json:"firstName" bson:"firstName"`
	LastName    string    `json:"lastName" bson:"lastName"`
	Password    string    `json:"password" bson:"password"`
	LastUpVote  time.Time `json:"lastUpVote" bson:"lastUpVote"`
}