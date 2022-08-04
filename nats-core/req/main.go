package main

import (
	"flag"
	"log"
	"time"

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
	goal := "POC NATS Core requestor"
	log.Println(goal)

	genMessages := flag.Int("msgs", 500, "message count")
	delay := flag.Int("d", 1, "delay per msgs")
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
	defer nc.Close()

	streamProvider := stream.NewNATSProvider(nc)
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	for i := 0; i < *genMessages; i++ {
		p := profitSvc.RequestProfit(subject, i)
		log.Printf("Requested on %s: %s\n", subject, p)
		time.Sleep(time.Duration(*delay) * time.Second)
	}

	return nil
}
