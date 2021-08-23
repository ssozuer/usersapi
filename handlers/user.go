package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"time"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}


func NewUserHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *UserHandler {
	return &UserHandler{
		collection:  collection,
		redisClient: redisClient,
		ctx:         ctx,
	}
}

// swagger:operation GET /users users listUsers
// Returns list of users
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func (handler *UserHandler) ListUsersHandler(c *gin.Context) {
	users := make([]models.User, 0)
	val, err := handler.redisClient.Get("users").Result()

	// If redis empty, touch the database
	if err == redis.Nil {
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}

		defer cur.Close(handler.ctx)

		for cur.Next(handler.ctx) {
			var user models.User
			cur.Decode(&user)
			users = append(users, user)
		}

		// Set users to Redis cache
		data, _ := json.Marshal(users)
		handler.redisClient.Set("users", string(data), 0)

		// Return users
		c.JSON(http.StatusOK, users)

	} else if err != nil {
		// Some error occurred
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		// Return users from Redis cache
		json.Unmarshal([]byte(val), &users)
		c.JSON(http.StatusOK, users)
	}
}

// swagger:operation POST /users users newUser
// Create a new user
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *UserHandler) NewUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Check if email exists

	h := sha256.New()

	user.ID = primitive.NewObjectID()
	user.Created = time.Now()
	user.Password = string(h.Sum([]byte(user.Password)))

	// Insert to database
	_, err := handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear Redis cache
	handler.redisClient.Del("users")

	// Return result
	c.JSON(http.StatusOK, user)
}

// swagger:operation DELETE /users/{id} users deleteUser
// Delete an existing user
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the user
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid user ID
func (handler *UserHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	deleteResult, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})

	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	if deleteResult.DeletedCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	// Return result
	c.JSON(http.StatusOK, gin.H{"message": "User has been deleted"})
}
