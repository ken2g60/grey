package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const TokenExpiration = 15 * time.Minute

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(user_id string, email string) (string, error) {

	token_lifespan := time.Now().Add(time.Minute * 30).Unix()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["email"] = email
	claims["exp"] = token_lifespan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))

}
