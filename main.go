package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Password    string    `json:"password"`
	LastUpVote  time.Time `json:"lastUpVote"`
}

type Post struct {
	CreatedAt time.Time `json:"createdAt"`
	UserID    string    `json:"userId"`
	FirstName string    `json:"firstName"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	UpVotes   []string  `json:"upVotes"`
}

type Comment struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	Content   string    `json:"content"`
}

func main() {
	app := fiber.New()

	// Connexion à la base de données MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Routes
	app.Post("/auth/register", register)
	app.Post("/auth/login", login)
	app.Get("/user/me", getUser)
	app.Put("/user/edit", updateUser)
	app.Delete("/user/remove", removeUser)
	app.Get("/post", getPosts)
	app.Post("/post", createPost)
	app.Get("/post/me", getUserPosts)
	app.Get("/post/:id", getPostByID)
	app.Delete("/post/:id", deletePost)
	app.Post("/post/vote/:id", votePost)
	app.Post("/comment", createComment)

	// Démarrage du serveur
	log.Fatal(app.Listen(":8080"))
}

func register(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func login(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func getUser(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func updateUser(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func removeUser(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func getPosts(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func createPost(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func getUserPosts(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func getPostByID(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func deletePost(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func votePost(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}

func createComment(c *fiber.Ctx) error {
	// À implémenter
	return c.SendStatus(http.StatusNotImplemented)
}
