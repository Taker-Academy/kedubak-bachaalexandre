package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"MyApi/encoding"
	"MyApi/database"
	"MyApi/models"
)

func RegisterUserHandler(c *fiber.Ctx) error {
	var newUser models.User

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Impossible de décoder les données de la requête",
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

func LoginUserHandler(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Impossible de décoder les données de la requête",
		})
	}
	user, err := database.GetUserByEmail(database.GetDB(), credentials.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erreur lors de la récupération de l'utilisateur depuis la base de données",
		})
	}
	if user == nil || !encoding.VerifyPassword(credentials.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Mauvaises identifiants",
		})
	}
	token, err := encoding.GenerateJWT(user.ID)
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
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
			},
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}