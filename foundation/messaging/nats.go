package messaging

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func Open(persistence bool, cfg Config, opt []nats.Option) (*nats.Conn, error) {
	var ServerURL string
	if persistence {
		ServerURL = cfg.JetstreamURL
	} else {
		ServerURL = cfg.ServerURL
	}
	nc, err := nats.Connect(ServerURL, opt...)
	return nc, err
}

func SetupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
