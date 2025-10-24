package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
	"time"
)

// Simple subscriber — receives every message
func SubscribeUserRegistered(js nats.JetStreamContext) {
	if js == nil {
		logger.Error("❌ JetStream context is nil")
		return
	}

	_, err := js.Subscribe("user.registered", func(m *nats.Msg) {
		logger.Info("📥 Received user.registered: " + string(m.Data))
		m.Ack()
	}, nats.Durable("user-reg-durable"), nats.DeliverAll(), nats.ManualAck())

	if err != nil {
		logger.Error("❌ Failed to subscribe: " + err.Error())
		return
	}

	logger.Info("✅ Subscribed to user.registered")
}

// Queue subscriber — only one instance in the group processes each message
func SubscribeQueue(js nats.JetStreamContext, group string) {
	if js == nil {
		logger.Error("❌ JetStream context is nil")
		return
	}

	_, err := js.QueueSubscribe("user.registered", group, func(m *nats.Msg) {
		logger.Info("📥 [" + group + "] received: " + string(m.Data))
		time.Sleep(1 * time.Second)
		m.Ack()
	}, nats.Durable("worker-durable-"+group), nats.ManualAck(), nats.AckWait(30*time.Second))

	if err != nil {
		logger.Error("❌ Failed to subscribe queue: " + err.Error())
		return
	}

	logger.Info("✅ Queue subscriber started for group: " + group)
}
