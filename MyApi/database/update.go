package database

import (
    "context"
    "errors"

    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

    "MyApi/models"
)

func UpdateUser(db *mongo.Database, userID string, updatedUser *models.User) error {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return errors.New("ID d'utilisateur invalide")
    }
    update := bson.M{
        "$set": bson.M{
            "email":     updatedUser.Email,
            "firstName": updatedUser.FirstName,
            "lastName":  updatedUser.LastName,
            "password":  updatedUser.Password,
        },
    }
    _, err = db.Collection("User").UpdateOne(
        context.TODO(),
        bson.M{"_id": objID},
        update,
    )
    if err != nil {
        return err
    }
    return nil
}