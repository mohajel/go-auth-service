package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
)

func Connect(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Error("❌ Failed to connect to NATS: " + err.Error())
		return nil
	}

	logger.Info("✅ Connected to NATS at " + url)
	return nc
}
