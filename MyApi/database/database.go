package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"MyApi/models"
)

func ConnectDb() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://alexandre:dbPassword@kedubak.kte5jfc.mongodb.net/?retryWrites=true&w=majority&appName=kedubak").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		client.Disconnect(context.TODO())
		return nil, err
	}

	userCollection := client.Database("kedubak").Collection("User")
	if _, err := userCollection.InsertOne(context.TODO(), models.User{}); err != nil {
		client.Disconnect(context.TODO())
		return nil, err
	}

	postCollection := client.Database("kedubak").Collection("Post")
	if _, err := postCollection.InsertOne(context.TODO(), models.Post{}); err != nil {
		client.Disconnect(context.TODO())
		return nil, err
	}

	fmt.Println("Connexion réussie à MongoDB et collections créées avec succès!")

	return client, nil
}