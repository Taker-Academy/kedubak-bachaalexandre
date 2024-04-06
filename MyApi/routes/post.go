package routes

import (
    "context"
	"time"
	"errors"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
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
			"ok": false,
            "error": "Erreur interne du serveur",
        })
    }
    defer cursor.Close(context.Background())
    
    var posts []models.Post
    for cursor.Next(context.Background()) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok": false,
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
			"ok": false,
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
			"ok": false,
            "error": "Erreur interne du serveur lors de la récupération des informations de l'utilisateur",
        })
    }
    var requestBody struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
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
			"ok": false,
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
			"ok": false,
            "error": "Erreur interne du serveur lors de la récupération des posts de l'utilisateur",
        })
    }
    defer cursor.Close(context.Background())

    var userPosts []models.Post
    for cursor.Next(context.Background()) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok": false,
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
			"ok": false,
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

func GetPostDetailsHandler(c *fiber.Ctx) error {
    postID := c.Params("id")
    if postID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }
    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, ID de l'élément invalide",
        })
    }
    var post models.Post
    db := database.GetDB()
    err = db.Collection("Post").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&post)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok": false,
                "error": "Élément non trouvé",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur",
        })
    }
	if post.Comments == nil {
		post.Comments = []models.Comment{}
	}
	if post.UpVotes == nil {
		post.UpVotes = []string{}
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
    return c.JSON(fiber.Map{
        "ok":   true,
        "data": responseData,
    })
}

func DeletePostHandler(c *fiber.Ctx) error {
    postID := c.Params("id")
    if postID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }
    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, ID de l'élément invalide",
        })
    }
    claims := c.Locals("user").(*encoding.CustomClaims)
    userID := claims.UserID
    db := database.GetDB()
    var post models.Post
    err = db.Collection("Post").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&post)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok": false,
                "error": "Élément non trouvé",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur",
        })
    }
    if post.UserID != userID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"ok": false,
            "error": "L'utilisateur n'est pas le propriétaire de l'élément",
        })
    }
    _, err = db.Collection("Post").DeleteOne(context.Background(), bson.M{"_id": objID})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur lors de la suppression de l'élément",
        })
    }
	if post.Comments == nil {
		post.Comments = []models.Comment{}
	}
	if post.UpVotes == nil {
		post.UpVotes = []string{}
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
        "removed":   true,
    }
    return c.JSON(fiber.Map{
        "ok":   true,
        "data": responseData,
    })
}
func VotePostHandler(c *fiber.Ctx) error {
    postID := c.Params("id")
    if postID == "" {
        return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"ok": false,
            "error": "ID invalide",
        })
    }
    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"ok": false,
            "error": "ID invalide",
        })
    }
    claims := c.Locals("user").(*encoding.CustomClaims)
    userID := claims.UserID
    db := database.GetDB()
    var post models.Post
    err = db.Collection("Post").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&post)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok": false,
                "error": "Élément non trouvé",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur",
        })
    }
    for _, v := range post.UpVotes {
        if v == userID {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "error": "Vous avez déjà voté pour ce post",
            })
        }
    }
    _, err = db.Collection("Post").UpdateOne(
        context.Background(),
        bson.M{"_id": objID},
        bson.M{"$push": bson.M{"upVotes": userID}},
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur lors de l'enregistrement du vote",
        })
    }
    return c.JSON(fiber.Map{
        "ok":      true,
        "message": "post upvoted",
    })
}

func CreateCommentHandler(c *fiber.Ctx) error {
    postID := c.Params("id")
    if postID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }
    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, ID de l'élément invalide",
        })
    }
    var requestBody struct {
        Content string `json:"content"`
    }
    if err := c.BodyParser(&requestBody); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, paramètres manquants ou invalides",
        })
    }
    if requestBody.Content == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok": false,
            "error": "Mauvaise requête, contenu du commentaire manquant",
        })
    }
    claims := c.Locals("user").(*encoding.CustomClaims)
    userID := claims.UserID
    user, err := database.GetUserByID(database.GetDB(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur lors de la récupération des informations de l'utilisateur",
        })
    }
    comment := models.Comment{
        CreatedAt: time.Now(),
        ID:        primitive.NewObjectID().Hex(),
        FirstName: user.FirstName,
        Content:   requestBody.Content,
    }
    db := database.GetDB()
    _, err = db.Collection("Post").UpdateOne(
        context.Background(),
        bson.M{"_id": objID},
        bson.M{"$push": bson.M{"comments": comment}},
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
            "error": "Erreur interne du serveur lors de la création du commentaire",
        })
    }
    responseData := map[string]interface{}{
        "firstName": comment.FirstName,
        "content":   comment.Content,
        "createdAt": comment.CreatedAt,
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "ok":   true,
        "data": responseData,
    })
}

