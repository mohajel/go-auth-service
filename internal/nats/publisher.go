package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
)


func PublishUserRegistered(nc *nats.Conn, userID string) {
	if nc == nil {
		logger.Error("NATS connection is nil, cannot publish user.registered event")
		return
	}

	err := nc.Publish("user.registered", []byte(userID))
	if err != nil {
		logger.Error("Failed to publish user.registered: " + err.Error())
		return
	}

	logger.Info("Published user.registered event for user: " + userID)
}

func PublishUserLoggedOut(nc *nats.Conn, userID string) {
	if nc == nil {
		logger.Error("NATS connection is nil, cannot publish user.logged_out event")
		return
	}

	err := nc.Publish("user.logged_out", []byte(userID))
	if err != nil {
		logger.Error("Failed to publish user.logged_out: " + err.Error())
		return
	}

	logger.Info("Published user.logged_out event for user: " + userID)
}
