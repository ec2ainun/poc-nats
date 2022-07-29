package main

import (
	"flag"
	"log"
	"runtime"

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
	// Exit the main goroutine but allow subscribe to continue running
	runtime.Goexit()
}

func run() error {
	goal := "POC NATS Core subscriber"
	log.Println(goal)

	flag.Parse()
	args := flag.Args()
	subject := args[0]

	conf, err := config.Read("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "err open config")
	}

	opts := []nats.Option{nats.Name(goal)}
	opts = messaging.SetupConnOptions(opts)

	nc, err := messaging.Open(false, conf, opts)
	if err != nil {
		return err
	}

	streamProvider := stream.NewNATSProvider(nc)
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	profitSvc.SubscribeProfit(subject)

	return nil
}
