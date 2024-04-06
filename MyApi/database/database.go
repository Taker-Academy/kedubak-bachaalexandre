package database

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "MyApi/models"
)

var db *mongo.Database

func ConnectDb() (*mongo.Database, error) {
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
    db = client.Database("kedubak")
    initCollections()
    return db, nil
}

func initCollections() {
    userCollection := db.Collection("User")
    userCollection.InsertOne(context.TODO(), models.User{})
    postCollection := db.Collection("Post")
    postCollection.InsertOne(context.TODO(), models.Post{})
}

func GetDB() *mongo.Database {
    return db
}
