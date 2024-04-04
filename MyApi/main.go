package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"MyApi/database"
	"MyApi/encoding"
	"MyApi/routes"
)

func setupRoutes(app *fiber.App) {
	app.Post("/auth/register", routes.RegisterUserHandler)
	app.Post("/auth/login", routes.LoginUserHandler)
	app.Get("/user/me", encoding.Authenticate, routes.GetUserInfoHandler)
	app.Put("/user/edit", encoding.Authenticate, routes.EditUserInfoHandler)
	app.Delete("/user/remove", encoding.Authenticate, routes.RemoveUserHandler)
}

func main() {
	database.ConnectDb()
	app := fiber.New()
	app.Use(cors.New())
	setupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
