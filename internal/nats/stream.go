package nats

import (
	"fmt"
	"time"

	"go-auth-service/pkg/logger"

	"github.com/nats-io/nats.go"
)

const (
	StreamName = "USER_EVENTS"
)

func EnsureStream(js nats.JetStreamContext) error {
	if js == nil {
		return fmt.Errorf("JetStream context is nil")
	}

	stream, err := js.StreamInfo(StreamName)
	if err == nil && stream != nil {
		logger.Info("Stream " + StreamName + " already exists")
		return nil
	}

	cfg := &nats.StreamConfig{
		Name:      StreamName,
		Subjects:  []string{"user.>"}, // wildcard for all user.* events
		Retention: nats.LimitsPolicy,
		Storage:   nats.FileStorage,
		MaxAge:    24 * time.Hour,
		MaxMsgs:   1000000,
	}

	stream, err = js.AddStream(cfg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	logger.Info("Stream " + StreamName + " created successfully")
	return nil
}
