package routes

import (
    "context"
	"time"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "MyApi/models"
	"MyApi/database"
	"MyApi/encoding"
)

func GetPostsHandler(c *fiber.Ctx) error {
    db := database.GetDB()
    opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
    cursor, err := db.Collection("Post").Find(context.Background(), bson.M{"userId": bson.M{"$ne": ""}}, opts)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur interne du serveur",
        })
    }
    defer cursor.Close(context.Background())
    
    var posts []models.Post
    for cursor.Next(context.Background()) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Erreur de décodage des posts",
            })
        }
        if post.Comments == nil {
            post.Comments = []models.Comment{}
        }
		if post.UpVotes == nil {
            post.UpVotes= []string{}
        }
        posts = append(posts, post)
    }
    if err := cursor.Err(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur de curseur lors de la récupération des posts",
        })
    }
    
    var responseData []map[string]interface{}
    if len(posts) == 0 {
        responseData = []map[string]interface{}{}
    } else {
        for _, post := range posts {
            postData := map[string]interface{}{
                "_id":       post.ID,
                "createdAt": post.CreatedAt,
                "userId":    post.UserID,
                "firstName": post.FirstName,
                "title":     post.Title,
                "content":   post.Content,
                "comments":  post.Comments,
                "upVotes":   post.UpVotes,
            }
            responseData = append(responseData, postData)
        }
    }
    response := fiber.Map{
        "ok":   true,
        "data": responseData,
    }
    return c.JSON(response)
}

func CreatePostHandler(c *fiber.Ctx) error {
    claims := c.Locals("user").(*encoding.CustomClaims)
    userID := claims.UserID
    user, err := database.GetUserByID(database.GetDB(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur interne du serveur lors de la récupération des informations de l'utilisateur",
        })
    }
    var requestBody struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }
    post := models.Post{
        CreatedAt: time.Now(),
        Title:     requestBody.Title,
        Content:   requestBody.Content,
        UserID:    user.ID,
        FirstName: user.FirstName,
        Comments:  []models.Comment{},
        UpVotes:   []string{},
    }
    db := database.GetDB()
    _, err = db.Collection("Post").InsertOne(context.Background(), post)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur interne du serveur lors de la création du post",
        })
    }
    responseData := map[string]interface{}{
        "_id":       post.ID,
        "createdAt": post.CreatedAt,
        "userId":    post.UserID,
        "firstName": post.FirstName,
        "title":     post.Title,
        "content":   post.Content,
        "comments":  post.Comments,
        "upVotes":   post.UpVotes,
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "ok":   true,
        "data": responseData,
    })
}

func GetUserPostsHandler(c *fiber.Ctx) error {
    claims := c.Locals("user").(*encoding.CustomClaims)
    userID := claims.UserID
    db := database.GetDB()
    opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
    cursor, err := db.Collection("Post").Find(context.Background(), bson.M{"userId": userID}, opts)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur interne du serveur lors de la récupération des posts de l'utilisateur",
        })
    }
    defer cursor.Close(context.Background())

    var userPosts []models.Post
    for cursor.Next(context.Background()) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Erreur de décodage des posts de l'utilisateur",
            })
        }
        if post.Comments == nil {
            post.Comments = []models.Comment{}
        }
		if post.UpVotes == nil {
            post.UpVotes = []string{}
        }
        userPosts = append(userPosts, post)
    }
    if err := cursor.Err(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Erreur de curseur lors de la récupération des posts de l'utilisateur",
        })
    }
    if len(userPosts) == 0 {
        return c.JSON(fiber.Map{
            "ok":   true,
            "data": []map[string]interface{}{},
        })
    }

    var responseData []map[string]interface{}
    for _, post := range userPosts {
        postData := map[string]interface{}{
            "_id":       post.ID,
            "createdAt": post.CreatedAt,
            "userId":    post.UserID,
            "firstName": post.FirstName,
            "title":     post.Title,
            "content":   post.Content,
            "comments":  post.Comments,
            "upVotes":   post.UpVotes,
        }
        responseData = append(responseData, postData)
    }

    return c.JSON(fiber.Map{
        "ok":   true,
        "data": responseData,
    })
}
