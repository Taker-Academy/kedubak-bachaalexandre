package routes

import (
	"time"

	"MyApi/database"
	"MyApi/encoding"
	"MyApi/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserHandler(c *fiber.Ctx) error {
	var newUser models.User

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Impossible de décoder les données de la requête",
		})
	}
	if newUser.Email == "" || newUser.Password == "" || newUser.FirstName == "" || newUser.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Les champs email, password, firstName et lastName sont obligatoires",
		})
	}
	newUser.CreatedAt = time.Now()
	hashedPassword, err := encoding.HashPassword(newUser.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erreur lors du chiffrement du mot de passe",
		})
	}
	newUser.Password = hashedPassword
	if err := database.SaveUser(newUser, "User"); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erreur lors de l'enregistrement de l'utilisateur dans la base de données",
		})
	}
	token, err := encoding.GenerateJWT(newUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erreur lors de la génération du token JWT",
		})
	}
	response := fiber.Map{
		"ok": true,
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"email":     newUser.Email,
				"firstName": newUser.FirstName,
				"lastName":  newUser.LastName,
			},
		},
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}
