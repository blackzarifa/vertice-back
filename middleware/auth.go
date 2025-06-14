package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthRequired creates a middleware that validates JWT tokens
// It extracts the token from Authorization header, validates it,
// and adds user info to the request context
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from "Authorization: Bearer <token>" header
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Missing authorization token"},
			)
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (any, error) {
				// Verify the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims and add to context for handlers to use
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("funcionario_id", claims["funcionario_id"])
			c.Set("nome", claims["nome"])
			c.Set("cargo", claims["cargo"])
		}

		c.Next()
	}
}

// extractToken gets the JWT token from Authorization header
func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
