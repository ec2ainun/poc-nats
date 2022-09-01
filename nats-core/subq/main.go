package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
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
	goal := "POC NATS Core subscriber Queue"
	log.Println(goal)

	var queueName = flag.String("q", "test-queue", "Queue Group Name")
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

	nc, err := messaging.Open(false, conf, opts)
	if err != nil {
		return err
	}

	streamProvider := stream.NewNATSQueueProvider(nc, *queueName)
	streamSvc := connector.NewStreamConnector(streamProvider)
	profitSvc := service.NewProfitService(streamSvc)

	var wg sync.WaitGroup
	go func() {
		<-ctx.Done()
		if err := nc.Drain(); err != nil {
			log.Fatalf("Error on drain: %v", err)
		}
		<-isNATSConnClosed
		<-time.After(5 * time.Second) //add leeway to gracefull shutdown
		wg.Done()
		os.Exit(0)
	}()

	wg.Add(1)
	go profitSvc.QueueSubscribeProfit(subject)
	wg.Wait()

	return nil
}
