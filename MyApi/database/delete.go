package database

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrUserNotFound = fmt.Errorf("utilisateur non trouvé")

func RemoveUser(db *mongo.Database, userID string) error {
    collection := db.Collection("user")
    filter := bson.M{"_id": userID}
    res, err := collection.DeleteOne(context.Background(), filter)
    if err != nil {
        return fmt.Errorf("erreur lors de la suppression de l'utilisateur: %v", err)
    }
    if res.DeletedCount == 0 {
        return ErrUserNotFound
    }
    return nil
}
