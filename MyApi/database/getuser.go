package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

func GetUserByID(db *mongo.Database, userID string) (*models.User, error) {
    var user models.User

    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }
	fmt.Println("objectID:", objectID)
    collection := db.Collection("User")
    filter := bson.M{"_id": objectID}
    err = collection.FindOne(context.Background(), filter).Decode(&user)

    if err == mongo.ErrNoDocuments {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    return &user, nil
}