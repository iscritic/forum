package session

import (
	"fmt"
	"forum/internal/repository"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

func CreateSession(db *repository.Storage, user *repository.User) (string, error) {

	u4, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %v", err)
	}
	log.Printf("generated Version 4 UUID %v", u4)

	var sess repository.Session

	sess.UserID = user.ID
	sess.SessionToken = u4

	expTime := time.Now().Add(20 * time.Minute)
	timeStamp := expTime.Format("2006-01-02 15:04:05")

	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStamp)
	if err != nil {
		return u4.String(), err
	}
	sess.ExpiredAt = parsedTime

	err = db.CreateSession(sess)
	if err != nil {
		return "", err
	}

	return u4.String(), nil
}
