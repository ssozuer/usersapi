package handlers

import (
	"context"
	"crypto/sha256"
	"net/http"
	"os"
	"time"
	"user-service/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var JWT_SECRET string = "JWT_SECRET"

type AuthHandler struct {
	collection   *mongo.Collection
	ctx 			   context.Context		 
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx: ctx,
	}
}

// swagger:operation POST /login login
// Logins the user
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *AuthHandler) LoginHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h := sha256.New()
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"email": user.Email,
		"password": string(h.Sum([]byte(user.Password))),
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password!"})
		return
	}

	tokenExpirationTime := time.Now().Add(time.Hour)
	claims := &models.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{ExpiresAt: tokenExpirationTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv(JWT_SECRET)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := models.JWTOutput{
		Token: tokenString,
		Expires: tokenExpirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}


func (handler *AuthHandler) AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv(JWT_SECRET)), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}