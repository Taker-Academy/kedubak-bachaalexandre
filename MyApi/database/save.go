package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"MyApi/models"
)

var (
	mongoURI     = "mongodb+srv://alexandre:dbPassword@kedubak.kte5jfc.mongodb.net/?retryWrites=true&w=majority&appName=kedubak"
	databaseName = "kedubak"
)

func SaveUser(user models.User, collectionName string) error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	collection := client.Database(databaseName).Collection(collectionName)
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}