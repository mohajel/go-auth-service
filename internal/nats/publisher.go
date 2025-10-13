package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
)

func PublishUserRegistered(js nats.JetStreamContext, userID string) {
	_, err := js.Publish("user.registered", []byte(userID))
	if err != nil {
		logger.Error("Failed to publish user.registered: " + err.Error())
		return
	}
	logger.Info("Published user.registered for user " + userID)
}

func PublishUserLoggedOut(js nats.JetStreamContext, userID string) {
	_, err := js.Publish("user.logged_out", []byte(userID))
	if err != nil {
		logger.Error("Failed to publish user.logged_out: " + err.Error())
		return
	}
	logger.Info("Published user.logged_out for user " + userID)
}
