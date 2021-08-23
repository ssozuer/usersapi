// User Service API
//
// This is a sample user service API.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//  Author: Selcuk Sozuer
//	Contact: <selcuk.sozuer@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"os"
	"user-service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Env. variables
var MONGO_URI string
var MONGO_DATABASE string
var USERS_COLLECTION string
var REDIS_URI string
var REDIS_PASSWORD string

// Handlers
var authHandler *handlers.AuthHandler
var userHandler *handlers.UserHandler

func init() {
	MONGO_URI        = os.Getenv("MONGO_URI")
	MONGO_DATABASE   = os.Getenv("MONGO_DATABASE")
	USERS_COLLECTION = os.Getenv("USERS_COLLECTION")
	REDIS_URI        = os.Getenv("REDIS_ADDRESS")
	REDIS_PASSWORD   = os.Getenv("REDIS_PASSWORD")

	// Connect to MongoDB Database
	ctx := context.Background()
	log.Println("Connecting to MongoDB Database...")
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(MONGO_URI),
	)

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connect to MongoDB Database Successfull")

	// User collections
	usersCollection := client.Database(MONGO_DATABASE).Collection(USERS_COLLECTION)

	// Connect to Redis
	log.Println("Connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr: REDIS_URI,
		Password: REDIS_PASSWORD,
		DB: 0,
	})
	status := redisClient.Ping()
	log.Println(status)

	// initialize handlers
	authHandler = handlers.NewAuthHandler(ctx, usersCollection)
	userHandler = handlers.NewUserHandler(ctx, usersCollection, redisClient)
}

func SetupServer() *gin.Engine {
	router := gin.Default()

	// Unauthorized API Endpoints
	router.GET("/users", userHandler.ListUsersHandler)
	router.POST("/users", userHandler.NewUserHandler)
	router.POST("/login", authHandler.LoginHandler)

	// Authorized API Endpoints
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleWare())
	{
		authorized.DELETE("/users/:id", userHandler.DeleteUserHandler)
	}
	return router
}

func main()  {
		SetupServer().Run(":8080")
}