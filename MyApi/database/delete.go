package database

import (
    "context"
    "errors"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = errors.New("utilisateur non trouvé")

func RemoveUser(db *mongo.Database, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID d'utilisateur invalide")
	}
	collection := db.Collection("User")
	filter := bson.M{"_id": objID}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}
	if res.DeletedCount == 0 {
		return fiber.NewError(fiber.StatusNotFound, ErrUserNotFound.Error())
	}

	return nil
}