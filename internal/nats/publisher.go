package nats

import (
	"fmt"
	"go-auth-service/pkg/logger"

	"github.com/nats-io/nats.go"
)

func PublishUserRegistered(js nats.JetStreamContext, userID string) error {
	if js == nil {
		return fmt.Errorf("JetStream context is nil")
	}

	_, err := js.Publish("user.registered", []byte(userID))
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
