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
	goal := "POC NATS Consumer Manager"
	log.Println(goal)

	cmd := flag.String("cmd", "", "command")
	sn := flag.String("sn", "", "stream name")
	cn := flag.String("cn", "", "consumer name")
	qn := flag.String("qn", "", "queue name")
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

	streamProvider := stream.NewJetStreamConsumerProvider(nc, jsctx, "", *sn, *cn, *qn)

	switch *cmd {
	case "add-push":
		if err := streamProvider.CreateConsumer("push", *sn, *cn, *subject, *qn); err != nil {
			log.Fatalf("err creating consumer: %s", err.Error())
		}
	case "add-pull":
		if err := streamProvider.CreateConsumer("pull", *sn, *cn, *subject, *qn); err != nil {
			log.Fatalf("err creating consumer: %s", err.Error())
		}
	case "delete":
		if err := streamProvider.DeleteConsumer(*sn, *cn); err != nil {
			log.Fatalf("err deleting consumer: %s", err.Error())
		}
	default:
		log.Fatalf("cmd not found, supported cmd: add, update, delete")
	}

	return nil
}
