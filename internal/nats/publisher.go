package nats

import (
	"encoding/json"
	"fmt"
	"go-auth-service/pkg/logger"

	"github.com/nats-io/nats.go"
)

type UserRegisteredEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name,omitempty"`
}

func PublishUserRegistered(js nats.JetStreamContext, userID string, email string, name string) error {
	if js == nil {
		return fmt.Errorf("JetStream context is nil")
	}

	event := UserRegisteredEvent{
		UserID: userID,
		Email:  email,
		Name:   name,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal user.registered event: " + err.Error())
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	_, err = js.Publish("user.registered", payload)
	if err != nil {
		logger.Error("Failed to publish user.registered: " + err.Error())
		return fmt.Errorf("failed to publish: %v", err)
	}

	logger.Info("Published user.registered for user " + userID)
	return nil
}

func PublishUserLoggedOut(js nats.JetStreamContext, userID string) {
	_, err := js.Publish("user.logged_out", []byte(userID))
	if err != nil {
		logger.Error("Failed to publish user.logged_out: " + err.Error())
		return
	}
	logger.Info("Published user.logged_out for user " + userID)
}
