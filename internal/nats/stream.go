package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
)

func EnsureStream(js nats.JetStreamContext) {
	if js == nil {
		logger.Error("JetStream context is nil")
		return
	}

	stream, err := js.StreamInfo("USER_EVENTS")
	if err == nil && stream != nil {
		logger.Info("Stream USER_EVENTS already exists")
		return
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "USER_EVENTS",
		Subjects:  []string{"user.registered", "user.logged_out"},
		Retention: nats.LimitsPolicy,
		Storage:   nats.FileStorage,
		MaxAge:    24 * time.Hour, 
	})

	if err != nil {
		logger.Error("Failed to create USER_EVENTS stream: " + err.Error())
		return
	}

	logger.Info("Stream USER_EVENTS created successfully")
	defer js.DeleteStream("USER_EVENTS")
}
