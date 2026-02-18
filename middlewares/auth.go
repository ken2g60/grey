package middlewares

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		bearerToken = strings.ReplaceAll(bearerToken, "Bearer ", "")

		if bearerToken == "" {
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Next()
			return
		}

		claimPayload, err := ValidateSessionToken(bearerToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("x-claim-payload", claimPayload)
		c.Set("x-token", bearerToken)

		c.Next()
	}
}

func ValidateSessionToken(tokenString string) (JwtSessionPayload, error) {

	var jwtKey = []byte(os.Getenv("SECRET_JWT"))
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return JwtSessionPayload{}, err
	}

	claims = token.Claims.(jwt.MapClaims)

	jsonString, _ := json.Marshal(claims)

	jwtPayload := JwtSessionPayload{}
	json.Unmarshal(jsonString, &jwtPayload)

	return jwtPayload, err

}
