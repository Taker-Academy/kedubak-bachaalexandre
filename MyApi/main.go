package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"MyApi/database"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	log.Fatal(app.Listen(":8080"))
}