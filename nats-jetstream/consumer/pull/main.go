package main

import (
	"flag"
	"log"

	connector "github.com/ec2ainun/poc-nats/business/connector"
	service "github.com/ec2ainun/poc-nats/business/service"
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
	goal := "POC NATS Jetstream consumer pull"
	log.Println(goal)

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

	streamProvider := stream.NewJetStreamConsumerProvider(nc, jsctx, "pull", *sn, *cn, *qn)
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	profitSvc.BatchSubscribeProfit(*subject, 3)

	return nil
}
