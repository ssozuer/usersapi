package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters users newUser
type User struct {
	ID 					primitive.ObjectID `json:"id"        bson:"_id"`
	FirstName   string  				   `json:"firstname" bson:"firstname"`
	LastName		string						 `json:"lastname"  bson:"lastname"`
	Email       string             `json:"email"     bson:"email"`
	Password    string             `json:"password"  bson:"password"`
	Created     time.Time          `json:"created"   bson:"created"`
}