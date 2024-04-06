package encoding

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok": false,
			"error": "Format de token JWT invalide",
		})
	}
	token := tokenParts[1]
	claims, err := VerifyJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok": false,
			"error": "mauvais token JWT",
		})
	}
	c.Locals("user", claims)
	return c.Next()
}

