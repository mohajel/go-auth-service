package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
)

func SubscribeUserRegistered(nc *nats.Conn) {
	if nc == nil {
		logger.Error("NATS connection is nil, cannot subscribe to user.registered")
		return
	}

	_, err := nc.Subscribe("user.registered", func(m *nats.Msg) {
		logger.Info("Received user.registered event: " + string(m.Data))
	})
	if err != nil {
		logger.Error("Failed to subscribe to user.registered: " + err.Error())
		return
	}

	logger.Info("Subscribed to user.registered event")
}
