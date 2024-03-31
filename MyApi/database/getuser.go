package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"MyApi/models"
)

func GetUserByEmail(db *mongo.Database, email string) (*models.User, error) {
	var user models.User

	collection := db.Collection("User")
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}