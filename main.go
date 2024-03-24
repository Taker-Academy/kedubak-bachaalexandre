package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
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

var client *mongo.Client
var ctx context.Context

func main() {
	app := fiber.New()

	// Connexion MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://alexandrebachaba:alex17MB@kedubak.2nmq2ix.mongodb.net/"))
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


func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// Fonction pour générer le token JWT
func generateJWT(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    })
    tokenString, err := token.SignedString([]byte("your_secret_key"))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// Fonction register mise à jour
func register(c *fiber.Ctx) error {
    // Parse request body
    type requestBody struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        FirstName string `json:"firstName"`
        LastName  string `json:"lastName"`
    }
    var body requestBody
    if err := c.BodyParser(&body); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "ok":    false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }

    // regerde si vide
    if body.Email == "" || body.Password == "" || body.FirstName == "" || body.LastName == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "ok":    false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }

    // check si il y a deja l'email
    existingUser := User{}
    err := client.Database("your_database_name").Collection("users").FindOne(ctx, bson.M{"email": body.Email}).Decode(&existingUser)
    if err == nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "ok":    false,
            "error": "Un utilisateur avec la même adresse e-mail existe déjà",
        })
    } else if err != mongo.ErrNoDocuments {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "ok":    false,
            "error": "Erreur interne du serveur",
        })
    }

    // on hash le passw
    hashedPassword, err := hashPassword(body.Password)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "ok":    false,
            "error": "Erreur interne du serveur",
        })
    }

    // sinon creation user
    newUser := User{
        ID:        primitive.NewObjectID().Hex(),
        CreatedAt: time.Now(),
        Email:     body.Email,
        FirstName: body.FirstName,
        LastName:  body.LastName,
        Password:  hashedPassword,
    }

    // insert dans la db
    _, err = client.Database("your_database_name").Collection("users").InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "ok":    false,
            "error": "Erreur interne du serveur",
        })
    }

    // Generation du token
    token, err := generateJWT(newUser.ID)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "ok":    false,
            "error": "Erreur interne du serveur",
        })
    }

    // Return ok
    return c.Status(http.StatusCreated).JSON(fiber.Map{
        "ok": true,
        "data": fiber.Map{
            "token": token,
            "user": fiber.Map{
                "email":     newUser.Email,
                "firstName": newUser.FirstName,
                "lastName":  newUser.LastName,
            },
        },
    })
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
