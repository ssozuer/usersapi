package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// swagger:parameters login
type Login struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Claims struct {
	Email  string `json:"email"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}