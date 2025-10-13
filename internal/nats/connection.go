package nats

import (
	"github.com/nats-io/nats.go"
	"go-auth-service/pkg/logger"
	"time"
)

func Connect(url string) *nats.Conn {
    nc, err := nats.Connect(url,
        nats.MaxReconnects(-1),
        nats.ReconnectWait(2 * time.Second),
    )
    if err != nil {
        logger.Error("NATS connection error: " + err.Error())
        return nil
    }
    logger.Info("Connected to NATS")
    return nc
}