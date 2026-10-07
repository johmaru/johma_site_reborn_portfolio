package models

import "time"

type CreateUserRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

type User struct {
	UID         string    `firestore:"uid"`
	Email       string    `firestore:"email"`
	DisplayName string    `firestore:"displayName"`
	CreatedAt   time.Time `firestore:"createdAt"`
}
