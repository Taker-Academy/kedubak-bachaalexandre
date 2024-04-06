package database

import (
    "context"
    "fmt"
    "errors"

    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = fmt.Errorf("utilisateur non trouvé")

func RemoveUser(db *mongo.Database, userID string) error {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return errors.New("ID d'utilisateur invalide")
    }
    collection := db.Collection("User")
    filter := bson.M{"_id": objID}
    res, err := collection.DeleteOne(context.Background(), filter)
    if err != nil {
        return fmt.Errorf("erreur1 lors de la suppression de l'utilisateur: %v", err)
    }
    if res.DeletedCount == 0 {
        return ErrUserNotFound
    }
    return nil
}
