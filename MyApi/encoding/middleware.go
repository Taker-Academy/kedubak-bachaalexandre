package encoding

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Format de token JWT invalide",
		})
	}
	token := tokenParts[1]
	fmt.Println("Token JWT:", token)
	claims, err := VerifyJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "mauvais token JWT",
		})
	}
	c.Locals("user", claims)
	return c.Next()
}

