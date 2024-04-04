package database

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/bson"

	"MyApi/models"
)

func UpdateUser(userID string, user *models.User) error {
    collection := db.Collection("User")
    filter := bson.M{"_id": userID}
    update := bson.M{}
    if user.FirstName != "" {
        update["firstName"] = user.FirstName
    }
    if user.LastName != "" {
        update["lastName"] = user.LastName
    }
    if user.Email != "" {
        update["email"] = user.Email
    }
    if user.Password != "" {
        update["password"] = user.Password
    }
    _, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": update})
    if err != nil {
        return fmt.Errorf("erreur lors de la mise à jour de l'utilisateur: %v", err)
    }

    return nil
}