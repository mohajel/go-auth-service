package nats_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	appnats "go-auth-service/internal/nats"
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
	_, err = js.Subscribe("user.registered", func(m *nats.Msg) {
		if string(m.Data) != "test-user" {
			t.Fatalf("unexpected message: %s", string(m.Data))
		}
		m.Ack()
		done <- true
	}, nats.Durable("test-durable"), nats.ManualAck())
	if err != nil {
		t.Fatal(err)
	}

	// Publish
	appnats.PublishUserRegistered(js, "test-user")

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}
