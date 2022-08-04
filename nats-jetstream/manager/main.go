package main

import (
	"flag"
	"log"

	config "github.com/ec2ainun/poc-nats/business/utils"
	stream "github.com/ec2ainun/poc-nats/connector/messaging"
	messaging "github.com/ec2ainun/poc-nats/foundation/messaging"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	goal := "POC NATS JetStream manager"
	log.Println(goal)

	cmd := flag.String("cmd", "", "command")
	streamName := flag.String("n", "", "name")
	subject := flag.String("s", "", "consumed subject")

	flag.Parse()

	conf, err := config.Read("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "err open config")
	}

	opts := []nats.Option{nats.Name(goal)}
	opts = messaging.SetupConnOptions(opts)

	nc, err := messaging.Open(true, conf, opts)
	if err != nil {
		return err
	}
	defer nc.Close()

	jsctx, err := nc.JetStream()
	if err != nil {
		return err
	}

	streamProvider := stream.NewJetStreamProvider(nc, jsctx)

	switch *cmd {
	case "add":
		streamProvider.Create(*streamName, *subject)
	case "update":
		streamProvider.Update(*streamName, *subject)
	case "delete":
		streamProvider.Delete(*streamName)
	default:
		log.Fatalf("cmd not found, supported cmd: add, update, delete")
	}

	return nil
}
