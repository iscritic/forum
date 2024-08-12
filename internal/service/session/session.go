package session

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

func CreateSession(db *repository.Storage, user *entity.User) (string, error) {
	// Delete all existing sessions for the user
	err := db.DeleteAllSessionsForUser(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to delete existing sessions: %v", err)
	}

	// Generate a new UUID for the session token
	u4, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %v", err)
	}
	log.Printf("generated Version 4 UUID %v", u4)

	// Create a new session struct
	var sess entity.Session
	sess.UserID = user.ID
	sess.SessionToken = u4

	// Set the expiration time for the session
	expTime := time.Now().Add(20 * time.Minute)
	sess.ExpiredAt = expTime

	// Store the session in the database
	err = db.CreateSession(sess)
	if err != nil {
		return "", err
	}

	return u4.String(), nil
}
