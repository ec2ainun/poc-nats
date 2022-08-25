package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
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
	goal := "POC NATS Jetstream publiser async"
	log.Println(goal)

	genMessages := flag.Int("msgs", 10, "message count")
	delay := flag.Int("d", 0, "delay per msgs")
	flag.Parse()
	args := flag.Args()
	subject := args[0]

	conf, err := config.Read("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "err open config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	isNATSConnClosed := make(chan bool)
	closed := func(nc *nats.Conn) {
		if nc.LastError() == nil {
			log.Printf("NATS conn closed")
			isNATSConnClosed <- true
		} else {
			log.Fatalf("Err clossing NATS conn: %v", nc.LastError())
		}
	}

	opts := []nats.Option{nats.Name(goal), nats.ClosedHandler(closed)}
	opts = messaging.SetupConnOptions(opts)

	nc, err := messaging.Open(true, conf, opts)
	if err != nil {
		return err
	}

	jsctx, err := nc.JetStream()
	if err != nil {
		return err
	}

	streamProvider := stream.NewJetStreamProducerProvider(nc, jsctx, "async")
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	go func() {
		<-ctx.Done()
		if err := nc.Drain(); err != nil {
			log.Fatalf("Error on drain: %v", err)
		}
		<-isNATSConnClosed
		os.Exit(0)
	}()

	for i := 0; i < *genMessages; i++ {
		profitSvc.PublishProfit(subject, i)
		time.Sleep(time.Duration(*delay) * time.Second)
	}

	return nil
}
