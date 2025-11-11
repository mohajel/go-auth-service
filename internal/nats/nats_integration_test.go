package nats_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	appnats "go-auth-service/internal/nats"

	"github.com/nats-io/nats.go"
)

func TestNATSConnectionAndPublish(t *testing.T) {
	nc := appnats.Connect("nats://localhost:4222")
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		t.Fatal(err)
	}
	appnats.EnsureStream(js)

	done := make(chan bool)

	// Subscriber
	durableName := fmt.Sprintf("test-durable-%d", time.Now().UnixNano())
	_, err = js.Subscribe("user.registered", func(m *nats.Msg) {
		var ev appnats.UserRegisteredEvent
		if err := json.Unmarshal(m.Data, &ev); err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}
		if ev.UserID != "test-user" {
			t.Fatalf("unexpected user id: %s", ev.UserID)
		}
		if ev.Email != "test@example.com" {
			t.Fatalf("unexpected email: %s", ev.Email)
		}
		m.Ack()
		done <- true
	}, nats.Durable(durableName), nats.ManualAck(), nats.DeliverNew())
	if err != nil {
		t.Fatal(err)
	}

	// Publish
	if err := appnats.PublishUserRegistered(js, "test-user", "test@example.com", "Test User"); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}
