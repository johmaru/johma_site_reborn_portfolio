package middleware

import (
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Auth *auth.Client
}

func NewAuthMiddleware(authClient *auth.Client) *AuthMiddleware {
	return &AuthMiddleware{Auth: authClient}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token, err := m.Auth.VerifyIDToken(ctx.Request.Context(), tokenString)
		if err != nil {
			log.Printf("Error verifying ID token: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			ctx.Abort()
			return
		}

		ctx.Set("user", token)
		ctx.Next()
	}
}

func (m *AuthMiddleware) OwnerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is missing"})
			c.Abort()
			return
		}

		firebaseToken := token.(*auth.Token)
		if firebaseToken.UID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not the owner"})
			c.Abort()
			return
		}

		c.Next()
	}
}
