package handlers

import (
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/johmaru/johma_site_reborn_portfolio/src/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Firestore *firestore.Client
	Auth      *auth.Client
}

func NewUserHandler(fs *firestore.Client, au *auth.Client) *UserHandler {
	return &UserHandler{
		Firestore: fs,
		Auth:      au,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.Auth.VerifyIDToken(ctx.Request.Context(), req.IDToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	uid := token.UID
	userRecord, err := h.Auth.GetUser(ctx.Request.Context(), uid)
	if err != nil {
		log.Printf("Error getting user record: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	newUser := models.User{
		UID:         userRecord.UID,
		Email:       userRecord.Email,
		DisplayName: userRecord.DisplayName,
		CreatedAt:   time.Now(),
	}

	_, err = h.Firestore.Collection("users").Doc(uid).Set(ctx.Request.Context(), newUser)
	if err != nil {
		log.Printf("Error creating user in Firestore: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "uid": uid})
}

// 本番環境では認証ミドルウェアで保護する
func (h *UserHandler) GetUser(c *gin.Context) {
	token, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	firebaseToken := token.(*auth.Token)
	userID := firebaseToken.UID

	requestedUserID := c.Param("id")
	if requestedUserID != "" && requestedUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not the owner"})
		return
	}

	ctx := c.Request.Context()

	doc, err := h.Firestore.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get user %s: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, doc.Data())
}
